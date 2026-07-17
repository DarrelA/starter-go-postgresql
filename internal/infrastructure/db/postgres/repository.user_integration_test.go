package postgres

import (
	"context"
	"errors"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/DarrelA/starter-go-postgresql/internal/domain/apperror"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

func TestUserRepositoryWithPostgres(t *testing.T) {
	dsn := os.Getenv("POSTGRES_TEST_DSN")
	if dsn == "" {
		t.Skip("POSTGRES_TEST_DSN is not set")
	}
	pool := newMigrationTestPool(t, dsn)
	if err := Migrate(context.Background(), pool); err != nil {
		t.Fatalf("migrate user repository schema: %v", err)
	}
	repository := NewUserRepository(pool)

	t.Run("missing users map to application error", func(t *testing.T) {
		if _, err := repository.GetByEmail(context.Background(), "missing@example.com"); !errors.Is(err, apperror.ErrUserNotFound) {
			t.Fatalf("expected missing-user error by email, got %v", err)
		}
		if _, err := repository.GetByUUID(context.Background(), uuid.New()); !errors.Is(err, apperror.ErrUserNotFound) {
			t.Fatalf("expected missing-user error by UUID, got %v", err)
		}
	})

	t.Run("PostgreSQL constraint errors remain inspectable", func(t *testing.T) {
		user := testDatabaseUser("oversized@example.com")
		user.FirstName = strings.Repeat("x", 256)
		err := repository.Save(context.Background(), user)
		var pgErr *pgconn.PgError
		if !errors.As(err, &pgErr) || pgErr.Code != "22001" {
			t.Fatalf("expected PostgreSQL value-too-long error, got %v", err)
		}
		if user.UUID != nil {
			t.Fatal("failed insert assigned a user UUID")
		}
	})

	t.Run("cancelled query returns cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := repository.GetByEmail(ctx, "cancelled@example.com")
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected context cancellation, got %v", err)
		}
	})

	t.Run("blocked insert respects deadline", func(t *testing.T) {
		lock, err := pool.Begin(context.Background())
		if err != nil {
			t.Fatalf("begin table-lock transaction: %v", err)
		}
		defer lock.Rollback(context.Background())
		if _, err := lock.Exec(context.Background(), "LOCK TABLE users IN ACCESS EXCLUSIVE MODE"); err != nil {
			t.Fatalf("lock users table: %v", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
		defer cancel()
		started := time.Now()
		err = repository.Save(ctx, testDatabaseUser("timeout@example.com"))
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("expected deadline error, got %v", err)
		}
		if elapsed := time.Since(started); elapsed > 2*time.Second {
			t.Fatalf("deadline took too long to cancel query: %s", elapsed)
		}
	})

	t.Run("concurrent duplicate registration creates one user", func(t *testing.T) {
		const workers = 8
		start := make(chan struct{})
		results := make(chan error, workers)
		var registrations sync.WaitGroup
		registrations.Add(workers)
		for range workers {
			go func() {
				defer registrations.Done()
				<-start
				results <- repository.Save(context.Background(), testDatabaseUser("race@example.com"))
			}()
		}
		close(start)
		registrations.Wait()
		close(results)

		succeeded, conflicted := 0, 0
		for err := range results {
			switch {
			case err == nil:
				succeeded++
			case errors.Is(err, apperror.ErrEmailConflict):
				conflicted++
			default:
				t.Fatalf("unexpected concurrent registration error: %v", err)
			}
		}
		if succeeded != 1 || conflicted != workers-1 {
			t.Fatalf("expected one success and %d conflicts, got %d and %d", workers-1, succeeded, conflicted)
		}
		var count int
		if err := pool.QueryRow(context.Background(),
			"SELECT count(*) FROM users WHERE email = $1", "race@example.com",
		).Scan(&count); err != nil {
			t.Fatalf("count concurrent user rows: %v", err)
		}
		if count != 1 {
			t.Fatalf("expected one concurrent user row, got %d", count)
		}
	})
}

func testDatabaseUser(email string) *entity.User {
	return &entity.User{FirstName: "Test", LastName: "User", Email: email, Password: "hash"}
}

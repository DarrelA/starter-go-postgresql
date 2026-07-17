package postgres

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/DarrelA/starter-go-postgresql/internal/domain/apperror"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestProviderIdentityRepositoryWithPostgres(t *testing.T) {
	dsn := os.Getenv("POSTGRES_TEST_DSN")
	if dsn == "" {
		t.Skip("POSTGRES_TEST_DSN is not set")
	}
	pool := newMigrationTestPool(t, dsn)
	if err := Migrate(context.Background(), pool); err != nil {
		t.Fatalf("migrate OAuth repository schema: %v", err)
	}
	repository := NewProviderIdentityRepository(pool)
	ctx := context.Background()

	var createdUUID uuid.UUID
	t.Run("new user", func(t *testing.T) {
		user, err := repository.FindLinkOrCreate(ctx, entity.ProviderIdentity{
			Provider: entity.OAuthProviderGoogle, Subject: "subject-new",
		}, &entity.User{FirstName: "New", LastName: "User", Email: "new@example.com"})
		if err != nil {
			t.Fatalf("create OAuth user: %v", err)
		}
		createdUUID = *user.UUID
		assertRowCount(t, pool, "users", 1)
		assertRowCount(t, pool, "provider_identities", 1)
	})

	t.Run("returning user", func(t *testing.T) {
		user, err := repository.FindLinkOrCreate(ctx, entity.ProviderIdentity{
			Provider: entity.OAuthProviderGoogle, Subject: "subject-new",
		}, &entity.User{FirstName: "Changed", LastName: "Profile", Email: "changed@example.com"})
		if err != nil {
			t.Fatalf("resolve returning OAuth user: %v", err)
		}
		if *user.UUID != createdUUID || user.Email != "new@example.com" {
			t.Fatalf("returning identity did not resolve stable local user: %#v", user)
		}
		assertRowCount(t, pool, "users", 1)
		assertRowCount(t, pool, "provider_identities", 1)
	})

	t.Run("link existing local user", func(t *testing.T) {
		var localUUID uuid.UUID
		if err := pool.QueryRow(ctx, queryInsertUser, "Local", "User", "local@example.com", "hash").Scan(&localUUID); err != nil {
			t.Fatalf("insert local user: %v", err)
		}
		user, err := repository.FindLinkOrCreate(ctx, entity.ProviderIdentity{
			Provider: entity.OAuthProviderGoogle, Subject: "subject-local",
		}, &entity.User{Email: "local@example.com"})
		if err != nil {
			t.Fatalf("link local user: %v", err)
		}
		if *user.UUID != localUUID {
			t.Fatalf("expected linked user %s, got %s", localUUID, user.UUID)
		}
	})

	t.Run("identity conflict", func(t *testing.T) {
		_, err := repository.FindLinkOrCreate(ctx, entity.ProviderIdentity{
			Provider: entity.OAuthProviderGoogle, Subject: "different-subject",
		}, &entity.User{Email: "local@example.com"})
		if !errors.Is(err, apperror.ErrOAuthIdentityConflict) {
			t.Fatalf("expected identity conflict, got %v", err)
		}
	})

	t.Run("rollback user when identity insert fails", func(t *testing.T) {
		usersBefore := tableRowCount(t, pool, "users")
		identitiesBefore := tableRowCount(t, pool, "provider_identities")
		_, err := repository.FindLinkOrCreate(ctx, entity.ProviderIdentity{
			Provider: strings.Repeat("p", 51), Subject: "rollback-subject",
		}, &entity.User{FirstName: "Rollback", LastName: "User", Email: "rollback@example.com"})
		if err == nil {
			t.Fatal("expected provider identity insert to fail")
		}
		if got := tableRowCount(t, pool, "users"); got != usersBefore {
			t.Fatalf("expected user insert rollback: before=%d after=%d", usersBefore, got)
		}
		if got := tableRowCount(t, pool, "provider_identities"); got != identitiesBefore {
			t.Fatalf("expected identity insert rollback: before=%d after=%d", identitiesBefore, got)
		}
	})

	t.Run("deadline during identity insert rolls back user", func(t *testing.T) {
		usersBefore := tableRowCount(t, pool, "users")
		identitiesBefore := tableRowCount(t, pool, "provider_identities")
		lock, err := pool.Begin(context.Background())
		if err != nil {
			t.Fatalf("begin identity-lock transaction: %v", err)
		}
		defer lock.Rollback(context.Background())
		if _, err := lock.Exec(context.Background(), "LOCK TABLE provider_identities IN SHARE MODE"); err != nil {
			t.Fatalf("lock provider identities: %v", err)
		}

		deadlineCtx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
		defer cancel()
		_, err = repository.FindLinkOrCreate(deadlineCtx, entity.ProviderIdentity{
			Provider: entity.OAuthProviderGoogle, Subject: "deadline-subject",
		}, &entity.User{FirstName: "Deadline", LastName: "User", Email: "deadline@example.com"})
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("expected OAuth transaction deadline, got %v", err)
		}
		if got := tableRowCount(t, pool, "users"); got != usersBefore {
			t.Fatalf("deadline did not roll back user: before=%d after=%d", usersBefore, got)
		}
		if got := tableRowCount(t, pool, "provider_identities"); got != identitiesBefore {
			t.Fatalf("deadline changed identities: before=%d after=%d", identitiesBefore, got)
		}
	})
}

func assertRowCount(t *testing.T, pool *pgxpool.Pool, table string, expected int) {
	t.Helper()
	if count := tableRowCount(t, pool, table); count != expected {
		t.Fatalf("expected %d %s row(s), got %d", expected, table, count)
	}
}

func tableRowCount(t *testing.T, pool *pgxpool.Pool, table string) int {
	t.Helper()
	var count int
	if err := pool.QueryRow(context.Background(), "SELECT count(*) FROM "+table).Scan(&count); err != nil {
		t.Fatalf("count %s: %v", table, err)
	}
	return count
}

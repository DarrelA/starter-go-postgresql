package postgres

import (
	"context"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestMigrateWithPostgres(t *testing.T) {
	dsn := os.Getenv("POSTGRES_TEST_DSN")
	if dsn == "" {
		t.Skip("POSTGRES_TEST_DSN is not set")
	}

	t.Run("clean database and repeated run", func(t *testing.T) {
		pool := newMigrationTestPool(t, dsn)
		if err := Migrate(context.Background(), pool); err != nil {
			t.Fatalf("first migration run: %v", err)
		}
		if err := Migrate(context.Background(), pool); err != nil {
			t.Fatalf("repeated migration run: %v", err)
		}
		assertMigrationState(t, pool)
	})

	t.Run("legacy schema adoption", func(t *testing.T) {
		pool := newMigrationTestPool(t, dsn)
		migrations, err := loadMigrations(migrationFiles)
		if err != nil {
			t.Fatalf("load migrations: %v", err)
		}
		if _, err := pool.Exec(context.Background(), migrations[0].sql); err != nil {
			t.Fatalf("create legacy schema: %v", err)
		}
		if err := Migrate(context.Background(), pool); err != nil {
			t.Fatalf("adopt legacy schema: %v", err)
		}
		assertMigrationState(t, pool)
	})

	t.Run("changed applied migration", func(t *testing.T) {
		pool := newMigrationTestPool(t, dsn)
		if err := Migrate(context.Background(), pool); err != nil {
			t.Fatalf("migrate: %v", err)
		}
		if _, err := pool.Exec(context.Background(),
			"UPDATE schema_migrations SET checksum = 'changed' WHERE version = 1",
		); err != nil {
			t.Fatalf("change checksum: %v", err)
		}
		if err := Migrate(context.Background(), pool); err == nil || !strings.Contains(err.Error(), "differs") {
			t.Fatalf("expected checksum mismatch, got %v", err)
		}
	})

	t.Run("concurrent migration jobs serialize", func(t *testing.T) {
		pool := newMigrationTestPool(t, dsn)
		const workers = 4
		errorsByWorker := make(chan error, workers)
		var jobs sync.WaitGroup
		jobs.Add(workers)
		for range workers {
			go func() {
				defer jobs.Done()
				errorsByWorker <- Migrate(context.Background(), pool)
			}()
		}
		jobs.Wait()
		close(errorsByWorker)
		for err := range errorsByWorker {
			if err != nil {
				t.Fatalf("concurrent migration: %v", err)
			}
		}
		assertMigrationState(t, pool)
	})

	t.Run("failed migration rolls back deployment transaction", func(t *testing.T) {
		pool := newMigrationTestPool(t, dsn)
		migrations := []migration{
			{version: 1, name: "000001_valid.up.sql", sql: "CREATE TABLE deployment_probe (id BIGINT PRIMARY KEY)", checksum: "valid"},
			{version: 2, name: "000002_invalid.up.sql", sql: "THIS IS NOT VALID SQL", checksum: "invalid"},
		}
		if err := applyMigrations(context.Background(), pool, migrations); err == nil {
			t.Fatal("expected failed deployment migration")
		}
		var tableName *string
		if err := pool.QueryRow(context.Background(), "SELECT to_regclass('deployment_probe')::text").Scan(&tableName); err != nil {
			t.Fatalf("inspect rolled-back migration: %v", err)
		}
		if tableName != nil {
			t.Fatalf("failed migration left table %q behind", *tableName)
		}
	})
}

func newMigrationTestPool(t *testing.T, dsn string) *pgxpool.Pool {
	t.Helper()
	ctx := context.Background()
	admin, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("open PostgreSQL admin pool: %v", err)
	}
	t.Cleanup(admin.Close)

	schema := "migration_test_" + strings.ReplaceAll(uuid.NewString(), "-", "")
	identifier := pgx.Identifier{schema}.Sanitize()
	if _, err := admin.Exec(ctx, "CREATE SCHEMA "+identifier); err != nil {
		t.Fatalf("create test schema: %v", err)
	}
	t.Cleanup(func() {
		if _, err := admin.Exec(context.Background(), "DROP SCHEMA "+identifier+" CASCADE"); err != nil {
			t.Errorf("drop test schema: %v", err)
		}
	})

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		t.Fatalf("parse PostgreSQL test DSN: %v", err)
	}
	config.ConnConfig.RuntimeParams["search_path"] = schema + ",public"
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		t.Fatalf("open schema-scoped PostgreSQL pool: %v", err)
	}
	t.Cleanup(pool.Close)
	return pool
}

func assertMigrationState(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	var migrationCount int
	if err := pool.QueryRow(context.Background(), "SELECT count(*) FROM schema_migrations").Scan(&migrationCount); err != nil {
		t.Fatalf("read migration ledger: %v", err)
	}
	if migrationCount != 3 {
		t.Fatalf("expected three applied migrations, got %d", migrationCount)
	}
	var usersTable string
	if err := pool.QueryRow(context.Background(), "SELECT to_regclass('users')::text").Scan(&usersTable); err != nil {
		t.Fatalf("find users table: %v", err)
	}
	if usersTable != "users" {
		t.Fatalf("expected users table, got %q", usersTable)
	}
	var identitiesTable string
	if err := pool.QueryRow(context.Background(), "SELECT to_regclass('provider_identities')::text").Scan(&identitiesTable); err != nil {
		t.Fatalf("find provider identities table: %v", err)
	}
	if identitiesTable != "provider_identities" {
		t.Fatalf("expected provider identities table, got %q", identitiesTable)
	}
}

package postgres

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func TestSchemaSupportsRepositoryQueriesWithPostgres(t *testing.T) {
	dsn := os.Getenv("POSTGRES_TEST_DSN")
	if dsn == "" {
		t.Skip("POSTGRES_TEST_DSN is not set")
	}
	pool := newMigrationTestPool(t, dsn)
	if err := Migrate(context.Background(), pool); err != nil {
		t.Fatalf("migrate schema contract: %v", err)
	}

	tx, err := pool.Begin(context.Background())
	if err != nil {
		t.Fatalf("begin query-plan transaction: %v", err)
	}
	defer tx.Rollback(context.Background())
	if _, err := tx.Exec(context.Background(), "SET LOCAL enable_seqscan = off"); err != nil {
		t.Fatalf("disable sequential scans for index assertions: %v", err)
	}
	fixtureUUID := seedQueryPlanFixtures(t, tx)

	assertPlanUses(t, tx, queryGetUser, []any{"plan-1@example.com"}, "users_email_key")
	assertPlanUses(t, tx, queryGetUserByID, []any{fixtureUUID}, "users_user_uuid_key")
	assertPlanUses(t, tx, queryUserByProviderIdentity,
		[]any{"google", "subject-1"}, "provider_identities_provider_subject_unique")
	assertPlanUses(t, tx, queryProviderSubjectByUser,
		[]any{"google", fixtureUUID}, "provider_identities_user_uuid_idx")

	assertColumnNotNullable(t, tx, "users", "user_uuid")
	assertConstraintExists(t, tx, "provider_identities_provider_not_blank")
	assertConstraintExists(t, tx, "provider_identities_subject_not_blank")
	assertConstraintExists(t, tx, "provider_identities_provider_user_unique")
	assertIndexExists(t, tx, "provider_identities_user_uuid_idx")
}

func seedQueryPlanFixtures(t *testing.T, tx pgx.Tx) uuid.UUID {
	t.Helper()
	_, err := tx.Exec(context.Background(), `
		INSERT INTO users(first_name, last_name, email, password)
		SELECT 'Plan', 'User', 'plan-' || value || '@example.com', 'hash'
		FROM generate_series(1, 200) AS value
	`)
	if err != nil {
		t.Fatalf("insert query-plan users: %v", err)
	}
	_, err = tx.Exec(context.Background(), `
		INSERT INTO provider_identities(user_uuid, provider, provider_subject)
		SELECT user_uuid, 'google',
			'subject-' || split_part(split_part(email, '@', 1), '-', 2)
		FROM users WHERE email LIKE 'plan-%@example.com'
	`)
	if err != nil {
		t.Fatalf("insert query-plan identities: %v", err)
	}
	if _, err := tx.Exec(context.Background(), "ANALYZE users; ANALYZE provider_identities"); err != nil {
		t.Fatalf("analyze query-plan fixtures: %v", err)
	}
	var userUUID uuid.UUID
	if err := tx.QueryRow(context.Background(),
		"SELECT user_uuid FROM users WHERE email = 'plan-1@example.com'",
	).Scan(&userUUID); err != nil {
		t.Fatalf("resolve query-plan user: %v", err)
	}
	return userUUID
}

func assertPlanUses(t *testing.T, tx pgx.Tx, query string, args []any, indexName string) {
	t.Helper()
	rows, err := tx.Query(context.Background(), "EXPLAIN (COSTS OFF) "+query, args...)
	if err != nil {
		t.Fatalf("explain query for %s: %v", indexName, err)
	}
	defer rows.Close()
	var plan strings.Builder
	for rows.Next() {
		var line string
		if err := rows.Scan(&line); err != nil {
			t.Fatalf("scan query plan for %s: %v", indexName, err)
		}
		plan.WriteString(line)
		plan.WriteByte('\n')
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("read query plan for %s: %v", indexName, err)
	}
	if !strings.Contains(plan.String(), indexName) {
		t.Fatalf("query plan did not use %s:\n%s", indexName, plan.String())
	}
}

func assertColumnNotNullable(t *testing.T, tx pgx.Tx, table, column string) {
	t.Helper()
	var nullable string
	err := tx.QueryRow(context.Background(), `
		SELECT is_nullable FROM information_schema.columns
		WHERE table_schema = current_schema() AND table_name = $1 AND column_name = $2
	`, table, column).Scan(&nullable)
	if err != nil {
		t.Fatalf("inspect %s.%s nullability: %v", table, column, err)
	}
	if nullable != "NO" {
		t.Fatalf("expected %s.%s to be NOT NULL, got %s", table, column, nullable)
	}
}

func assertConstraintExists(t *testing.T, tx pgx.Tx, name string) {
	t.Helper()
	var exists bool
	err := tx.QueryRow(context.Background(), `
		SELECT EXISTS (
			SELECT 1 FROM pg_constraint
			WHERE connamespace = current_schema()::regnamespace AND conname = $1
		)
	`, name).Scan(&exists)
	if err != nil {
		t.Fatalf("inspect constraint %s: %v", name, err)
	}
	if !exists {
		t.Fatalf("expected constraint %s", name)
	}
}

func assertIndexExists(t *testing.T, tx pgx.Tx, name string) {
	t.Helper()
	var exists bool
	err := tx.QueryRow(context.Background(), `
		SELECT EXISTS (
			SELECT 1 FROM pg_indexes
			WHERE schemaname = current_schema() AND indexname = $1
		)
	`, name).Scan(&exists)
	if err != nil {
		t.Fatalf("inspect index %s: %v", name, err)
	}
	if !exists {
		t.Fatalf("expected index %s", name)
	}
}

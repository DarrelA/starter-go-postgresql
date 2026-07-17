package postgres

import (
	"testing"
	"testing/fstest"
)

func TestEmbeddedMigrationsAreValid(t *testing.T) {
	migrations, err := loadMigrations(migrationFiles)
	if err != nil {
		t.Fatalf("load embedded migrations: %v", err)
	}
	if len(migrations) != 3 || migrations[0].version != 1 || migrations[1].version != 2 ||
		migrations[2].version != 3 {
		t.Fatalf("unexpected migrations: %#v", migrations)
	}
}

func TestLoadMigrationsSortsVersions(t *testing.T) {
	files := fstest.MapFS{
		"migrations/000002_second.up.sql": {Data: []byte("SELECT 2")},
		"migrations/000001_first.up.sql":  {Data: []byte("SELECT 1")},
	}

	migrations, err := loadMigrations(files)
	if err != nil {
		t.Fatalf("load migrations: %v", err)
	}
	if migrations[0].version != 1 || migrations[1].version != 2 {
		t.Fatalf("migrations are not sorted: %#v", migrations)
	}
}

func TestLoadMigrationsRejectsDuplicateVersions(t *testing.T) {
	files := fstest.MapFS{
		"migrations/000001_first.up.sql":  {Data: []byte("SELECT 1")},
		"migrations/000001_second.up.sql": {Data: []byte("SELECT 2")},
	}

	if _, err := loadMigrations(files); err == nil {
		t.Fatal("expected duplicate version error")
	}
}

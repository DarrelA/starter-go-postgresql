package main

import (
	"context"
	"errors"
	"testing"

	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
)

func TestRunReturnsConfigurationFailureBeforeExecution(t *testing.T) {
	configErr := errors.New("invalid database configuration")
	executed := false
	err := run(context.Background(), func() (*entity.PostgresDBConfig, error) {
		return nil, configErr
	}, func(context.Context, *entity.PostgresDBConfig) error {
		executed = true
		return nil
	})
	if !errors.Is(err, configErr) {
		t.Fatalf("expected configuration failure, got %v", err)
	}
	if executed {
		t.Fatal("migration executed after configuration failure")
	}
}

func TestRunReturnsMigrationFailure(t *testing.T) {
	migrationErr := errors.New("migration failed")
	err := run(context.Background(), func() (*entity.PostgresDBConfig, error) {
		return &entity.PostgresDBConfig{}, nil
	}, func(context.Context, *entity.PostgresDBConfig) error {
		return migrationErr
	})
	if !errors.Is(err, migrationErr) {
		t.Fatalf("expected migration failure, got %v", err)
	}
}

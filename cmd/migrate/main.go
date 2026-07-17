// Command migrate applies pending PostgreSQL schema migrations and exits.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	"github.com/DarrelA/starter-go-postgresql/internal/infrastructure/config"
	"github.com/DarrelA/starter-go-postgresql/internal/infrastructure/db/postgres"
	"github.com/rs/zerolog/log"
)

type configLoader func() (*entity.PostgresDBConfig, error)
type migrationExecutor func(context.Context, *entity.PostgresDBConfig) error

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := run(ctx, config.LoadPostgresConfig, executeMigrations); err != nil {
		log.Fatal().Err(err).Msg("migration job failed")
	}
	log.Info().Msg("database migrations are current")
}

func run(ctx context.Context, load configLoader, execute migrationExecutor) error {
	postgresConfig, err := load()
	if err != nil {
		return err
	}
	if err := execute(ctx, postgresConfig); err != nil {
		return fmt.Errorf("execute migrations: %w", err)
	}
	return nil
}

func executeMigrations(ctx context.Context, postgresConfig *entity.PostgresDBConfig) error {
	pool, err := postgres.Open(ctx, postgresConfig)
	if err != nil {
		return err
	}
	defer pool.Close()
	return postgres.Migrate(ctx, pool)
}

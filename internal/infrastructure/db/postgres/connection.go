// coverage:ignore file
// Testing with integration test
package postgres

import (
	"context"
	"fmt"

	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

const (
	errMsgCreatingConnectionPool = "unable to create connection pool"
	errMsgGreetingQuery          = "dbpool.QueryRow failed"
)

/*
dbpool is the database connection pool.
Package pgxpool is a concurrency-safe connection pool for pgx.
pgxpool implements a nearly identical interface to pgx connections.

- The `PostgresDB` is stateful because it holds a connection to the database (`pgxpool.Pool`). This dependency is injected into the repository to manage database operations.
- This pattern is useful for managing resources that have a lifecycle, like database connections.
*/
type PostgresDB struct {
	PostgresDBConfig *entity.PostgresDBConfig
	Dbpool           *pgxpool.Pool
}

func (p *PostgresDB) ConnectToPostgres(postgresDBConfig *entity.PostgresDBConfig) repository.RDBMS {
	connString := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s pool_max_conns=%s",
		postgresDBConfig.Username, postgresDBConfig.Password,
		postgresDBConfig.Host, postgresDBConfig.Port,
		postgresDBConfig.Name, postgresDBConfig.SslMode,
		postgresDBConfig.PoolMaxConns,
	)

	var err error
	dbpool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Error().Err(err).Msg(errMsgCreatingConnectionPool)
		panic(err)
	}

	var greeting string
	err = dbpool.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		log.Error().Err(err).Msg(errMsgGreetingQuery)
		panic(err)
	}

	log.Info().Msg("successfully connected to the Postgres database")
	return &PostgresDB{PostgresDBConfig: postgresDBConfig, Dbpool: dbpool}
}

func (p *PostgresDB) DisconnectFromPostgres() {
	if p.Dbpool != nil {
		p.Dbpool.Close()
		log.Info().Msg("PostgreSQL database connection closed")
	}
}

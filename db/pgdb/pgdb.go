package db

import (
	"context"
	"fmt"

	"github.com/DarrelA/starter-go-postgresql/configs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

/*
Embedding

Db is the database connection pool.
Package pgxpool is a concurrency-safe connection pool for pgx.
pgxpool implements a nearly identical interface to pgx connections.
*/
type PostgresDB struct {
	configs.PostgresDBConfig
}

var Dbpool *pgxpool.Pool

// NewDB creates a new PostgresDB instance with loaded config
func NewDB() *PostgresDB {
	return &PostgresDB{
		PostgresDBConfig: configs.PGDB,
	}
}

func (db *PostgresDB) Connect() {
	log.Info().
		Str("username", db.Username).
		Str("password", db.Password).
		Str("host", db.Host).
		Str("port", db.Port).
		Str("dbname", db.Name).
		Str("sslmode", db.SslMode).
		Str("pool_max_conns", db.PoolMaxConns).
		Msg("Database connection details")

	connString := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s pool_max_conns=%s",
		db.Username, db.Password, db.Host, db.Port, db.Name, db.SslMode, db.PoolMaxConns,
	)

	var err error
	Dbpool, err = pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Panic().Msg("unable to connect to the database: " + err.Error())
		panic(err)
	}

	log.Info().Msg("successfully connected to the Postgres database")
}

func (db *PostgresDB) Disconnect() {
	if Dbpool != nil {
		Dbpool.Close()
		log.Info().Msg("PostgreSQL database connection closed")
	}
}

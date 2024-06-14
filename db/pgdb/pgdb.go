package db

import (
	"context"
	"fmt"

	"github.com/DarrelA/starter-go-postgresql/configs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

// Db is the database connection pool.
// Package pgxpool is a concurrency-safe connection pool for pgx.
// pgxpool implements a nearly identical interface to pgx connections.
var Dbpool *pgxpool.Pool

func ConnectPostgres() {
	dbCfg := configs.PGDB

	connString := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s pool_max_conns=%s",
		dbCfg.Username, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.Name, dbCfg.SslMode, dbCfg.PoolMaxConns,
	)

	var err error
	Dbpool, err = pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Panic().Msg("unable to connect to the database: " + err.Error())
		panic(err)
	}

	log.Info().Msg("successfully connected to the Postgres database")
}

func DisconnectPostgres() {
	if Dbpool != nil {
		Dbpool.Close()
		log.Info().Msg("PostgreSQL database connection closed")
	}
}

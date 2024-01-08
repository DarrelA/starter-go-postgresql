package pgdb

import (
	"context"
	"fmt"
	"log"

	"github.com/DarrelA/starter-go-postgresql/configs"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Db is the database connection pool.
// Package pgxpool is a concurrency-safe connection pool for pgx.
// pgxpool implements a nearly identical interface to pgx connections.
var Db *pgxpool.Pool

func ConnectDatabase() {
	dbCfg := configs.PGDB

	connString := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s pool_max_conns=%s",
		dbCfg.Username, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.Name, dbCfg.SslMode, dbCfg.PoolMaxConns,
	)

	var err error
	Db, err = pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Printf("Unable to connect to the database: %v\n", err)
		panic(err)
	}
	log.Println("Successfully connected to the Postgres database!")
}

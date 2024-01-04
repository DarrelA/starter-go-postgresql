package pgdb

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// Db is the database connection pool.
// Package pgxpool is a concurrency-safe connection pool for pgx.
// pgxpool implements a nearly identical interface to pgx connections.
var Db *pgxpool.Pool

func ConnectDatabase() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error is occurred  on .env file please check")
	}

	user := os.Getenv("PGDB_USERNAME")
	pass := os.Getenv("PGDB_PASSWORD")
	host := os.Getenv("PGDB_HOST")
	port := os.Getenv("PGDB_PORT")
	dbname := os.Getenv("PGDB_NAME")
	sslmode := os.Getenv("PGDB_SSLMODE")
	pool_max_conns := os.Getenv("PGDB_POOL_MAX_CONNS")

	connString := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s pool_max_conns=%s",
		user, pass, host, port, dbname, sslmode, pool_max_conns,
	)

	Db, err = pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Printf("Unable to connect to the database: %v\n", err)
		panic(err)
	}
	log.Println("Successfully connected to the Postgres database!")
}

package pgdb

import (
	"context"
	"fmt"
	"log"

	"github.com/DarrelA/starter-go-postgresql/configs"
	"github.com/golang-migrate/migrate"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Db is the database connection pool.
// Package pgxpool is a concurrency-safe connection pool for pgx.
// pgxpool implements a nearly identical interface to pgx connections.
var Db *pgxpool.Pool

func ConnectDatabase() {
	dbCfg := configs.PGDB

	// Run migrations before establishing the connection pool
	if err := runMigrations(dbCfg); err != nil {
		log.Fatalf("Failed to run migrations: %v\n", err)
	}

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

func runMigrations(dbCfg configs.DBConfig) error {
	migrationDbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbCfg.Username, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.Name, dbCfg.SslMode)

	m, err := migrate.New("../migration", migrationDbURL)
	if err != nil {
		return err
	}

	// Apply all up migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

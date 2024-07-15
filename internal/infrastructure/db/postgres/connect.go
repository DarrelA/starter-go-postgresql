package postgres

import (
	"context"
	"fmt"

	"github.com/DarrelA/starter-go-postgresql/configs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

// @TODO: Try to use method receiver with config interface in domain service.
func NewRDBMS(db *configs.PostgresDBConfig) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s pool_max_conns=%s",
		db.Username, db.Password, db.Host, db.Port, db.Name, db.SslMode, db.PoolMaxConns,
	)

	var err error
	dbpool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Panic().Err(err).Msg("unable to create connection pool")
		panic(err)
	}

	var greeting string
	err = dbpool.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		log.Panic().Err(err).Msg("QueryRow failed")
		panic(err)
	}

	log.Info().Msg("successfully connected to the Postgres database")
	return dbpool, nil
}

func Disconnect(dbpool *pgxpool.Pool) {
	if dbpool != nil {
		dbpool.Close()
		log.Info().Msg("PostgreSQL database connection closed")
	}
}

package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/DarrelA/starter-go-postgresql/configs"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/repository"
	"github.com/DarrelA/starter-go-postgresql/internal/utils/err_rest"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

func (r *PostgresUserRepository) Disconnect() {
	if r.DbPool != nil {
		r.DbPool.Close()
		log.Info().Msg("PostgreSQL database connection closed")
	}
}

/*
Embedding

DbPool is the database connection pool.
Package pgxpool is a concurrency-safe connection pool for pgx.
pgxpool implements a nearly identical interface to pgx connections.
*/
type PostgresUserRepository struct {
	DbPool *pgxpool.Pool
}

func NewPostgresUserRepository(dbpool *pgxpool.Pool) repository.UserRepository {
	return &PostgresUserRepository{DbPool: dbpool}
}

var (
	queryInsertUser  = "INSERT INTO users(first_name, last_name, email, password) VALUES ($1, $2, $3, $4) RETURNING user_uuid;"
	queryGetUser     = "SELECT user_uuid, first_name, last_name, email, password FROM users WHERE email=$1;"
	queryGetUserByID = "SELECT user_uuid, first_name, last_name, email FROM users WHERE user_uuid=$1;"
)

// Create a method of the `User` type
func (r *PostgresUserRepository) SaveUser(user *entity.User) *err_rest.RestErr {
	var lastInsertUuid uuid.UUID
	err := r.DbPool.QueryRow(context.Background(), queryInsertUser, user.FirstName, user.LastName, user.Email, user.Password).Scan(&lastInsertUuid)

	if err != nil {
		var pgErr *pgconn.PgError

		// Check if `err` can be cast to `*pgconn.PgError` using `errors.As` before
		// attempting to access any fields of `pgErr`
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return err_rest.NewBadRequestError("email is already taken")
			}
		}

		log.Error().Err(err).Msg("pgdb_error")
		return err_rest.NewInternalServerError("something went wrong")
	}

	user.UUID = &lastInsertUuid
	return nil
}

func (r *PostgresUserRepository) GetUserByEmail(user *entity.User) *err_rest.RestErr {
	err := r.DbPool.QueryRow(context.Background(), queryGetUser, user.Email).
		Scan(&user.UUID, &user.FirstName, &user.LastName, &user.Email, &user.Password)

	if err != nil {
		if err == pgx.ErrNoRows {
			return err_rest.NewBadRequestError("the account has not been registered")
		}

		log.Error().Err(err).Msg("pgdb_error")
		return err_rest.NewInternalServerError("something went wrong")
	}

	return nil
}

func (r *PostgresUserRepository) GetUserByUUID(user *entity.User) *err_rest.RestErr {
	result := r.DbPool.QueryRow(context.Background(), queryGetUserByID, user.UUID)
	if err := result.Scan(&user.UUID, &user.FirstName, &user.LastName, &user.Email); err != nil {
		log.Error().Err(err).Msg("pgdb_error")
		return err_rest.NewInternalServerError("something went wrong")
	}

	return nil
}

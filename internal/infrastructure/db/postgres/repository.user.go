package postgres

import (
	"context"
	"errors"

	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	r "github.com/DarrelA/starter-go-postgresql/internal/domain/repository/postgres"
	"github.com/DarrelA/starter-go-postgresql/internal/utils/err_rest"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type PostgresUserRepository struct {
	dbpool *pgxpool.Pool
}

func NewUserRepository(dbpool *pgxpool.Pool) r.PostgresUserRepository {
	return &PostgresUserRepository{dbpool}
}

var (
	queryInsertUser  = "INSERT INTO users(first_name, last_name, email, password) VALUES ($1, $2, $3, $4) RETURNING user_uuid;"
	queryGetUser     = "SELECT user_uuid, first_name, last_name, email, password FROM users WHERE email=$1;"
	queryGetUserByID = "SELECT user_uuid, first_name, last_name, email FROM users WHERE user_uuid=$1;"
)

// Create a method of the `User` type
func (ur PostgresUserRepository) SaveUser(user *entity.User) *err_rest.RestErr {
	var lastInsertUuid uuid.UUID
	err := ur.dbpool.QueryRow(context.Background(), queryInsertUser, user.FirstName, user.LastName, user.Email, user.Password).Scan(&lastInsertUuid)

	if err != nil {
		var pgErr *pgconn.PgError

		// Check if `err` can be cast to `*pgconn.PgError` using `errors.As` before
		// attempting to access any fields of `pgErr`
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return err_rest.NewBadRequestError(err_rest.ErrMsgEmailIsAlreadyTaken)
			}
		}

		log.Error().Err(err).Msg("pgdb_error")
		return err_rest.NewInternalServerError(err_rest.ErrMsgSomethingWentWrong)
	}

	user.UUID = &lastInsertUuid
	return nil
}

func (ur PostgresUserRepository) GetUserByEmail(user *entity.User) *err_rest.RestErr {
	err := ur.dbpool.QueryRow(context.Background(), queryGetUser, user.Email).
		Scan(&user.UUID, &user.FirstName, &user.LastName, &user.Email, &user.Password)

	if err != nil {
		if err == pgx.ErrNoRows {
			return err_rest.NewBadRequestError("the account has not been registered")
		}

		log.Error().Err(err).Msg("pgdb_error")
		return err_rest.NewInternalServerError(err_rest.ErrMsgSomethingWentWrong)
	}

	return nil
}

func (ur PostgresUserRepository) GetUserByUUID(user *entity.User) *err_rest.RestErr {
	result := ur.dbpool.QueryRow(context.Background(), queryGetUserByID, user.UUID)
	if err := result.Scan(&user.UUID, &user.FirstName, &user.LastName, &user.Email); err != nil {
		log.Error().Err(err).Msg("pgdb_error")
		return err_rest.NewInternalServerError(err_rest.ErrMsgSomethingWentWrong)
	}

	return nil
}

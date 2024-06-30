package users

import (
	"context"
	"errors"

	pgdb "github.com/DarrelA/starter-go-postgresql/db/pgdb"
	"github.com/DarrelA/starter-go-postgresql/internal/utils/err_rest"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog/log"
)

var (
	queryInsertUser  = "INSERT INTO users(first_name, last_name, email, password) VALUES ($1, $2, $3, $4) RETURNING user_uuid;"
	queryGetUser     = "SELECT user_uuid, first_name, last_name, email, password FROM users WHERE email=$1;"
	queryGetUserByID = "SELECT user_uuid, first_name, last_name, email FROM users WHERE user_uuid=$1;"
)

/*
This function saves a new user to the database. The user's data is provided
via a pointer to the `User` struct. The data flow typically starts from
`users_controller.go`, goes through `users_service.go`, and finally reaches
`users_dao.go` where it interacts with the database.
*/
func (user *User) Save() *err_rest.RestErr {
	var lastInsertUuid uuid.UUID
	err := pgdb.Dbpool.QueryRow(context.Background(), queryInsertUser, user.FirstName, user.LastName, user.Email, user.Password).Scan(&lastInsertUuid)

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

func (user *User) GetByEmail() *err_rest.RestErr {
	err := pgdb.Dbpool.QueryRow(context.Background(), queryGetUser, user.Email).
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

func (user *User) GetByUUID() *err_rest.RestErr {
	result := pgdb.Dbpool.QueryRow(context.Background(), queryGetUserByID, user.UUID)
	if err := result.Scan(&user.UUID, &user.FirstName, &user.LastName, &user.Email); err != nil {
		log.Error().Err(err).Msg("pgdb_error")
		return err_rest.NewInternalServerError("something went wrong")
	}

	return nil
}

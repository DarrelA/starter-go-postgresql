package users

import (
	"context"

	pgdb "github.com/DarrelA/starter-go-postgresql/db/pgdb"
	"github.com/DarrelA/starter-go-postgresql/internal/utils/errors"
	"github.com/google/uuid"
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
func (user *User) Save() *errors.RestErr {
	var lastInsertUuid *uuid.UUID
	err := pgdb.Dbpool.QueryRow(context.Background(), queryInsertUser, user.FirstName, user.LastName, user.Email, user.Password).Scan(&lastInsertUuid)
	if err != nil {
		log.Error().Msg("pgdb_error: " + err.Error())
		return errors.NewInternalServerError("something went wrong")
	}

	user.UUID = lastInsertUuid
	return nil
}

func (user *User) GetByEmail() *errors.RestErr {
	result := pgdb.Dbpool.QueryRow(context.Background(), queryGetUser, user.Email)
	if err := result.Scan(&user.UUID, &user.FirstName, &user.LastName, &user.Email, &user.Password); err != nil {
		log.Error().Msg("pgdb_error: " + err.Error())
		return errors.NewInternalServerError("something went wrong")
	}

	return nil
}

func (user *User) GetByUUID() *errors.RestErr {
	result := pgdb.Dbpool.QueryRow(context.Background(), queryGetUserByID, user.UUID)
	if err := result.Scan(&user.UUID, &user.FirstName, &user.LastName, &user.Email); err != nil {
		log.Error().Msg("pgdb_error: " + err.Error())
		return errors.NewInternalServerError("something went wrong")
	}

	return nil
}

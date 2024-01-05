package users

import (
	"context"

	"github.com/DarrelA/starter-go-postgresql/db/pgdb"
	"github.com/DarrelA/starter-go-postgresql/internal/utils/errors"
)

var (
	queryInsertUser = "INSERT INTO users(first_name, last_name, email, password) VALUES ($1, $2, $3, $4) RETURNING id;"
	queryGetUser    = "SELECT id, first_name, last_name, email, password FROM users WHERE email=$1;"
)

/*
This function saves a new user to the database. The user's data is provided
via a pointer to the `User` struct. The data flow typically starts from
`users_controller.go`, goes through `users_service.go`, and finally reaches
`users_dao.go` where it interacts with the database.
*/
func (user *User) Save() *errors.RestErr {
	var lastInsertId int64
	err := pgdb.Db.QueryRow(context.Background(), queryInsertUser, user.FirstName, user.LastName, user.Email, user.Password).Scan(&lastInsertId)
	if err != nil {
		return errors.NewInternalServerError("database error: " + err.Error())
	}

	user.ID = lastInsertId
	return nil
}

func (user *User) GetByEmail() *errors.RestErr {
	result := pgdb.Db.QueryRow(context.Background(), queryGetUser, user.Email)
	if getErr := result.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password); getErr != nil {
		return errors.NewInternalServerError("database error")
	}

	return nil
}

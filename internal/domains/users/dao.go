package users

import (
	"context"
	"log"

	"github.com/DarrelA/starter-go-postgresql/db/pgdb"
	"github.com/DarrelA/starter-go-postgresql/internal/utils/errors"
)

var (
	queryInsertUser = "INSERT INTO users(first_name, last_name, email, password) VALUES ($1, $2, $3, $4) RETURNING id;"
)

func (user *User) Save() *errors.RestErr {
	// values come from pointer to the `User`
	// `users_controller.go` to `users_service.go` to `users_dao.go`

	var lastInsertId int64
	err := pgdb.Db.QueryRow(context.Background(), queryInsertUser, user.FirstName, user.LastName, user.Email, user.Password).Scan(&lastInsertId)
	if err != nil {
		return errors.NewInternalServerError("database error: " + err.Error())
	}

	log.Println(lastInsertId)

	// user.ID = lastInsertId
	return nil
}

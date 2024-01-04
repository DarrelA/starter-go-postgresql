package services

import (
	"github.com/DarrelA/starter-go-postgresql/internal/domains/users"
	"github.com/DarrelA/starter-go-postgresql/internal/utils/errors"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(user users.User) (*users.User, *errors.RestErr) {
	if err := user.Validate(); err != nil {
		return nil, err
	}

	pwSlice, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return nil, errors.NewBadRequestError(("failed to encrypt the password"))
	}

	// parse from byte to string
	user.Password = string(pwSlice[:])

	if err := user.Save(); err != nil {
		return nil, err
	}

	return &user, nil
}

/*
@TODO:
	- Rename `services.go` to `auth.go`
	- Move `token.go` from `utils` folder to `services` folder
*/

package services

import (
	"github.com/DarrelA/starter-go-postgresql/internal/domains/users"
	"github.com/DarrelA/starter-go-postgresql/internal/utils/errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(payload users.RegisterInput) (*users.UserResponse, *errors.RestErr) {
	newUser := &users.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  payload.Password,
	}

	if err := newUser.Validate(); err != nil {
		return nil, err
	}

	pwSlice, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 14)
	if err != nil {
		return nil, errors.NewBadRequestError(("failed to encrypt the password"))
	}

	// parse from byte to string
	newUser.Password = string(pwSlice[:])

	if err := newUser.Save(); err != nil {
		return nil, err
	}

	userResponse := &users.UserResponse{
		UUID:      newUser.UUID,
		FirstName: newUser.FirstName,
		LastName:  newUser.LastName,
		Email:     newUser.Email,
	}

	return userResponse, nil
}

func GetUser(user users.LoginInput) (*users.UserResponse, *errors.RestErr) {
	result := &users.User{Email: user.Email}
	if err := result.GetByEmail(); err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password)); err != nil {
		return nil, errors.NewBadRequestError("failed to decrypt the password")
	}

	userResponse := &users.UserResponse{
		UUID:      result.UUID,
		FirstName: result.FirstName,
		LastName:  result.LastName,
		Email:     result.Email,
	}

	return userResponse, nil
}

func GetUserByUUID(userUuid string) (*users.User, *errors.RestErr) {
	uuidPointer, err := uuid.Parse(userUuid)
	if err != nil {
		return nil, errors.NewInternalServerError("Error parsing UUID: " + err.Error())
	}

	result := &users.User{UUID: &uuidPointer}

	if err := result.GetByUUID(); err != nil {
		return nil, err
	}
	return result, nil
}

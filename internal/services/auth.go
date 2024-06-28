package services

import (
	"github.com/DarrelA/starter-go-postgresql/internal/domains/users"
	"github.com/DarrelA/starter-go-postgresql/internal/utils/err_rest"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(payload users.RegisterInput) (*users.UserResponse, *err_rest.RestErr) {
	if err := users.ValidateStruct(payload); err != nil {
		return nil, err
	}

	newUser := &users.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  payload.Password,
	}

	pwSlice, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 14)
	if err != nil {
		log.Error().Err(err).Msg("bcrypt_error")
		return nil, err_rest.NewInternalServerError(("something went wrong"))
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

func GetUser(user users.LoginInput) (*users.UserResponse, *err_rest.RestErr) {
	result := &users.User{Email: user.Email}
	if err := result.GetByEmail(); err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password)); err != nil {
		return nil, err_rest.NewBadRequestError("invalid credentials")
	}

	userResponse := &users.UserResponse{
		UUID:      result.UUID,
		FirstName: result.FirstName,
		LastName:  result.LastName,
		Email:     result.Email,
	}

	return userResponse, nil
}

func GetUserByUUID(userUuid string) (*users.User, *err_rest.RestErr) {
	uuidPointer, err := uuid.Parse(userUuid)
	if err != nil {
		log.Error().Err(err).Msg("uuid_error")
		return nil, err_rest.NewUnprocessableEntityError(("something went wrong"))
	}

	result := &users.User{UUID: &uuidPointer}

	if err := result.GetByUUID(); err != nil {
		return nil, err
	}
	return result, nil
}

package factory

import (
	user "github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/repository"
	"github.com/DarrelA/starter-go-postgresql/internal/utils/err_rest"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type UserFactory struct {
	ur repository.UserRepository
}

func NewUserFactory(ur repository.UserRepository) *UserFactory {
	return &UserFactory{ur}
}

func (uf *UserFactory) CreateUser(payload user.RegisterInput) (*user.UserResponse, *err_rest.RestErr) {
	newUser := &user.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  payload.Password,
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 14)
	if err != nil {
		log.Error().Err(err).Msg("bcrypt_error")
		return nil, err_rest.NewInternalServerError("something went wrong")
	}

	newUser.Password = string(hashedPassword)

	if err := uf.ur.SaveUser(newUser); err != nil {
		return nil, err
	}

	userResponse := &user.UserResponse{
		UUID:      newUser.UUID,
		FirstName: newUser.FirstName,
		LastName:  newUser.LastName,
		Email:     newUser.Email,
	}

	return userResponse, nil
}

func (uf *UserFactory) GetUser(u user.LoginInput) (*user.UserResponse, *err_rest.RestErr) {
	result := &user.User{Email: u.Email}
	if err := uf.ur.GetUserByEmail(result); err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(u.Password)); err != nil {
		return nil, err_rest.NewBadRequestError("invalid credentials")
	}

	userResponse := &user.UserResponse{
		UUID:      result.UUID,
		FirstName: result.FirstName,
		LastName:  result.LastName,
		Email:     result.Email,
	}

	return userResponse, nil
}

func (uf *UserFactory) GetUserByUUID(userUuid string) (*user.User, *err_rest.RestErr) {
	uuidPointer, err := uuid.Parse(userUuid)
	if err != nil {
		log.Error().Err(err).Msg("uuid_error")
		return nil, err_rest.NewUnprocessableEntityError(("something went wrong"))
	}

	result := &user.User{UUID: &uuidPointer}

	if err := uf.ur.GetUserByUUID(result); err != nil {
		return nil, err
	}
	return result, nil
}

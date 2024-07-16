package factory

import (
	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/factory"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/repository"
	dto "github.com/DarrelA/starter-go-postgresql/internal/interface/transport/dto"
	"github.com/DarrelA/starter-go-postgresql/internal/utils/err_rest"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

/*
The `UserFactory` is responsible for creating and retrieving `User` entities.
The factory pattern here is used to encapsulate the creation logic,
including hashing passwords and calling the repository to save or fetch users.

The factory logic is part of the domain layer as it encapsulates domain-specific creation and retrieval logic.
This is appropriate for the domain layer, as it deals with core business logic.
*/
type UserFactory struct {
	JWTConfig *entity.JWTConfig
	ur        repository.UserRepository
}

/*
The factory interacts with the `UserRepository` interface to perform persistence operations.
This adheres to the principle of dependency inversion,
where the factory depends on an abstraction rather than a concrete implementation.
*/
func NewUserFactory(JWTConfig *entity.JWTConfig, ur repository.UserRepository) factory.UserFactory {
	return &UserFactory{JWTConfig, ur}
}

func (uf *UserFactory) GetJWTConfig() *entity.JWTConfig {
	return uf.JWTConfig
}

func (uf *UserFactory) CreateUser(payload dto.RegisterInput) (*dto.UserResponse, *err_rest.RestErr) {
	newUser := &entity.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  payload.Password,
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 14)
	if err != nil {
		log.Error().Err(err).Msg("bcrypt_error")
		return nil, err_rest.NewInternalServerError(err_rest.ErrMsgSomethingWentWrong)
	}

	newUser.Password = string(hashedPassword)

	if err := uf.ur.SaveUser(newUser); err != nil {
		return nil, err
	}

	userResponse := &dto.UserResponse{
		UUID:      newUser.UUID,
		FirstName: newUser.FirstName,
		LastName:  newUser.LastName,
		Email:     newUser.Email,
	}

	return userResponse, nil
}

func (uf *UserFactory) GetUser(u dto.LoginInput) (*dto.UserResponse, *err_rest.RestErr) {
	result := &entity.User{Email: u.Email}
	if err := uf.ur.GetUserByEmail(result); err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(u.Password)); err != nil {
		return nil, err_rest.NewBadRequestError(err_rest.ErrMsgInvalidCredentials)
	}

	userResponse := &dto.UserResponse{
		UUID:      result.UUID,
		FirstName: result.FirstName,
		LastName:  result.LastName,
		Email:     result.Email,
	}

	return userResponse, nil
}

func (uf *UserFactory) GetUserByUUID(userUuid string) (*entity.User, *err_rest.RestErr) {
	uuidPointer, err := uuid.Parse(userUuid)
	if err != nil {
		log.Error().Err(err).Msg("uuid_error")
		return nil, err_rest.NewUnprocessableEntityError((err_rest.ErrMsgSomethingWentWrong))
	}

	result := &entity.User{UUID: &uuidPointer}

	if err := uf.ur.GetUserByUUID(result); err != nil {
		return nil, err
	}
	return result, nil
}

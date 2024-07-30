package factory

import (
	dto "github.com/DarrelA/starter-go-postgresql/internal/application/dto"
	"github.com/DarrelA/starter-go-postgresql/internal/application/factory"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	errConst "github.com/DarrelA/starter-go-postgresql/internal/domain/error"
	restDomainErr "github.com/DarrelA/starter-go-postgresql/internal/domain/error/transport/http"
	repo "github.com/DarrelA/starter-go-postgresql/internal/domain/repository/postgres"
	password "github.com/DarrelA/starter-go-postgresql/internal/infrastructure/bcrypt"
	restInterfaceErr "github.com/DarrelA/starter-go-postgresql/internal/interface/transport/http/error"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
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
	ur        repo.PostgresUserRepository
}

/*
The factory interacts with the `PostgresUserRepository` interface to perform persistence operations.
This adheres to the principle of dependency inversion,
where the factory depends on an abstraction rather than a concrete implementation.
*/
func NewUserFactory(JWTConfig *entity.JWTConfig, ur repo.PostgresUserRepository) factory.UserFactory {
	return &UserFactory{JWTConfig, ur}
}

func (uf *UserFactory) GetJWTConfig() *entity.JWTConfig {
	return uf.JWTConfig
}

func (uf *UserFactory) CreateUser(payload dto.RegisterInput) (*dto.UserResponse, *restDomainErr.RestErr) {
	newUser := &entity.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  payload.Password,
	}

	hashedPassword, err := password.HashPassword(newUser.Password)
	if err != nil {
		return nil, restInterfaceErr.NewInternalServerError(errConst.ErrMsgSomethingWentWrong)
	}

	newUser.Password = hashedPassword

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

func (uf *UserFactory) GetUserByEmail(u dto.LoginInput) (*dto.UserResponse, *restDomainErr.RestErr) {
	result := &entity.User{Email: u.Email}
	if err := uf.ur.GetUserByEmail(result); err != nil {
		return nil, err
	}

	if err := password.VerifyPassword(result.Password, u.Password); err != nil {
		return nil, restInterfaceErr.NewBadRequestError(err.Error())
	}

	userResponse := &dto.UserResponse{
		UUID:      result.UUID,
		FirstName: result.FirstName,
		LastName:  result.LastName,
		Email:     result.Email,
	}

	return userResponse, nil
}

func (uf *UserFactory) GetUserByUUID(userUuid string) (*entity.User, *restDomainErr.RestErr) {
	uuidPointer, err := uuid.Parse(userUuid)
	if err != nil {
		log.Error().Err(err).Msg(errConst.ErrUUIDError)
		return nil, restInterfaceErr.NewUnprocessableEntityError((errConst.ErrMsgSomethingWentWrong))
	}

	result := &entity.User{UUID: &uuidPointer}

	if err := uf.ur.GetUserByUUID(result); err != nil {
		return nil, err
	}
	return result, nil
}

// coverage:ignore file
// Testing with integration test
package service

import (
	dto "github.com/DarrelA/starter-go-postgresql/internal/application/dto"
	appSvc "github.com/DarrelA/starter-go-postgresql/internal/application/service"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	errConst "github.com/DarrelA/starter-go-postgresql/internal/domain/error"
	restDomainErr "github.com/DarrelA/starter-go-postgresql/internal/domain/error/transport/http"
	repo "github.com/DarrelA/starter-go-postgresql/internal/domain/repository/postgres"
	password "github.com/DarrelA/starter-go-postgresql/internal/infrastructure/bcrypt"
	restInterfaceErr "github.com/DarrelA/starter-go-postgresql/internal/interface/transport/http/error"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type UserService struct {
	JWTConfig *entity.JWTConfig
	ur        repo.PostgresUserRepository
}

func NewUserService(JWTConfig *entity.JWTConfig, ur repo.PostgresUserRepository) appSvc.UserService {
	return &UserService{JWTConfig, ur}
}

func (us *UserService) GetJWTConfig() *entity.JWTConfig {
	return us.JWTConfig
}

func (us *UserService) CreateUser(payload dto.RegisterInput) (*dto.UserResponse, *restDomainErr.RestErr) {
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

	if err := us.ur.SaveUser(newUser); err != nil {
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

func (us *UserService) GetUserByEmail(u dto.LoginInput) (*dto.UserResponse, *restDomainErr.RestErr) {
	result := &entity.User{Email: u.Email}
	if err := us.ur.GetUserByEmail(result); err != nil {
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

func (us *UserService) GetUserByUUID(userUuid string) (*entity.User, *restDomainErr.RestErr) {
	uuidPointer, err := uuid.Parse(userUuid)
	if err != nil {
		log.Error().Err(err).Msg(errConst.ErrUUIDError)
		return nil, restInterfaceErr.NewUnprocessableEntityError((errConst.ErrMsgSomethingWentWrong))
	}

	result := &entity.User{UUID: &uuidPointer}

	if err := us.ur.GetUserByUUID(result); err != nil {
		return nil, err
	}
	return result, nil
}

package service

import (
	dto "github.com/DarrelA/starter-go-postgresql/internal/application/dto"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	restDomainErr "github.com/DarrelA/starter-go-postgresql/internal/domain/error/transport/http"
)

type UserService interface {
	GetJWTConfig() *entity.JWTConfig
	CreateUser(payload dto.RegisterInput) (*dto.UserResponse, *restDomainErr.RestErr)
	GetUserByEmail(u dto.LoginInput) (*dto.UserResponse, *restDomainErr.RestErr)
	GetUserByUUID(userUuid string) (*entity.User, *restDomainErr.RestErr)
}

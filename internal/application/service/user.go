package service

import (
	dto "github.com/DarrelA/starter-go-postgresql/internal/application/dto"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	restErr "github.com/DarrelA/starter-go-postgresql/internal/error"
)

type UserService interface {
	GetJWTConfig() *entity.JWTConfig
	CreateUser(payload dto.RegisterInput) (*dto.UserResponse, *restErr.RestErr)
	GetUserByEmail(u dto.LoginInput) (*dto.UserResponse, *restErr.RestErr)
	GetUserByUUID(userUuid string) (*entity.User, *restErr.RestErr)
}

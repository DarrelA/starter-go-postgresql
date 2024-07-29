package factory

import (
	dto "github.com/DarrelA/starter-go-postgresql/internal/application/repository/dto"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	restDomainErr "github.com/DarrelA/starter-go-postgresql/internal/domain/error/transport/http"
)

type UserFactory interface {
	GetJWTConfig() *entity.JWTConfig
	CreateUser(payload dto.RegisterInput) (*dto.UserResponse, *restDomainErr.RestErr)
	GetUserByEmail(u dto.LoginInput) (*dto.UserResponse, *restDomainErr.RestErr)
	GetUserByUUID(userUuid string) (*entity.User, *restDomainErr.RestErr)
}

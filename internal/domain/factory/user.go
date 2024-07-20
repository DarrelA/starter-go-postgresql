package factory

import (
	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	restDomainErr "github.com/DarrelA/starter-go-postgresql/internal/domain/error/transport/http"
	dto "github.com/DarrelA/starter-go-postgresql/internal/interface/transport/dto"
)

type UserFactory interface {
	GetJWTConfig() *entity.JWTConfig
	CreateUser(payload dto.RegisterInput) (*dto.UserResponse, *restDomainErr.RestErr)
	GetUser(u dto.LoginInput) (*dto.UserResponse, *restDomainErr.RestErr)
	GetUserByUUID(userUuid string) (*entity.User, *restDomainErr.RestErr)
}

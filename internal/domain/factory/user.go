package factory

import (
	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	dto "github.com/DarrelA/starter-go-postgresql/internal/interface/transport/dto"
	"github.com/DarrelA/starter-go-postgresql/internal/utils/err_rest"
)

type UserFactory interface {
	GetJWTConfig() *entity.JWTConfig
	CreateUser(payload dto.RegisterInput) (*dto.UserResponse, *err_rest.RestErr)
	GetUser(u dto.LoginInput) (*dto.UserResponse, *err_rest.RestErr)
	GetUserByUUID(userUuid string) (*entity.User, *err_rest.RestErr)
}

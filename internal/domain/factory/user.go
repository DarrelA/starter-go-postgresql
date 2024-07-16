package factory

import (
	user "github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	dto "github.com/DarrelA/starter-go-postgresql/internal/interface/transport/dto"
	"github.com/DarrelA/starter-go-postgresql/internal/utils/err_rest"
)

type UserFactory interface {
	CreateUser(payload dto.RegisterInput) (*dto.UserResponse, *err_rest.RestErr)
	GetUser(u dto.LoginInput) (*dto.UserResponse, *err_rest.RestErr)
	GetUserByUUID(userUuid string) (*user.User, *err_rest.RestErr)
}

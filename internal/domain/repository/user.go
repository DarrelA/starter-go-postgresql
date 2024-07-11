package repository

import (
	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	"github.com/DarrelA/starter-go-postgresql/internal/utils/err_rest"
)

type UserRepository interface {
	SaveUser(user *entity.User) *err_rest.RestErr
	GetUserByEmail(user *entity.User) *err_rest.RestErr
	GetUserByUUID(user *entity.User) *err_rest.RestErr
}

package repository

import (
	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	"github.com/DarrelA/starter-go-postgresql/internal/utils/err_rest"
)

/*
The `UserRepository` interface defines the contract for the repository layer.

Defining the repository interface in the domain layer is appropriate as
it represents a boundary for persistence operations
*/
type UserRepository interface {
	SaveUser(user *entity.User) *err_rest.RestErr
	GetUserByEmail(user *entity.User) *err_rest.RestErr
	GetUserByUUID(user *entity.User) *err_rest.RestErr
}

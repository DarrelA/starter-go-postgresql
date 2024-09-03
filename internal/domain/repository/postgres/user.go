package repository

import (
	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	restErr "github.com/DarrelA/starter-go-postgresql/internal/error"
)

/*
The `UserRepository` interface defines the contract for the repository layer.

Defining the repository interface in the domain layer is appropriate as
it represents a boundary for persistence operations
*/
type PostgresUserRepository interface {
	SaveUser(user *entity.User) *restErr.RestErr
	GetUserByEmail(user *entity.User) *restErr.RestErr
	GetUserByUUID(user *entity.User) *restErr.RestErr
}

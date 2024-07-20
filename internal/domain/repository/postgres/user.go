package repository

import (
	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	restDomainErr "github.com/DarrelA/starter-go-postgresql/internal/domain/error/transport/http"
)

/*
The `UserRepository` interface defines the contract for the repository layer.

Defining the repository interface in the domain layer is appropriate as
it represents a boundary for persistence operations
*/
type PostgresUserRepository interface {
	SaveUser(user *entity.User) *restDomainErr.RestErr
	GetUserByEmail(user *entity.User) *restDomainErr.RestErr
	GetUserByUUID(user *entity.User) *restDomainErr.RestErr
}

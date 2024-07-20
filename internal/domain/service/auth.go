package service

import (
	"time"

	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	restDomainErr "github.com/DarrelA/starter-go-postgresql/internal/domain/error/transport/http"
)

/*
The `TokenService` interface define the contract for authentication-related operations.
*/
type TokenService interface {
	CreateToken(userUuid string, ttl time.Duration, privateKey string) (*entity.Token, *restDomainErr.RestErr)
	ValidateToken(token string, publicKey string) (*entity.Token, *restDomainErr.RestErr)
}

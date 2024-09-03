package service

import (
	"time"

	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	restErr "github.com/DarrelA/starter-go-postgresql/internal/error"
)

/*
The `TokenService` interface define the contract for authentication-related operations.
*/
type TokenService interface {
	CreateToken(userUuid string, ttl time.Duration, privateKey string) (*entity.Token, *restErr.RestErr)
	ValidateToken(token string, publicKey string) (*entity.Token, *restErr.RestErr)
}

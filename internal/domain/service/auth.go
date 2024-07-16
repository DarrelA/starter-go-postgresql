package service

import (
	"time"

	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	"github.com/DarrelA/starter-go-postgresql/internal/utils/err_rest"
)

/*
The `TokenService` interface define the contract for authentication-related operations.
*/
type TokenService interface {
	CreateToken(user_uuid string, ttl time.Duration, privateKey string) (*entity.Token, *err_rest.RestErr)
	ValidateToken(token string, publicKey string) (*entity.Token, *err_rest.RestErr)
}

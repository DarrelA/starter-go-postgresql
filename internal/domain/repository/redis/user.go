package repository

import restErr "github.com/DarrelA/starter-go-postgresql/internal/error"

type RedisUserRepository interface {
	SetUserUUID(tokenUUID string, userUUID string, expiresIn int64) *restErr.RestErr
	GetUserUUID(tokenUUID string) (string, *restErr.RestErr)
	DelUserUUID(tokenUUID string, accessTokenUUID string) (int64, *restErr.RestErr)
}

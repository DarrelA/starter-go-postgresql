package repository

import restDomainErr "github.com/DarrelA/starter-go-postgresql/internal/domain/error/transport/http"

type RedisUserRepository interface {
	SetUserUUID(tokenUUID string, userUUID string, expiresIn int64) *restDomainErr.RestErr
	GetUserUUID(tokenUUID string) (string, *restDomainErr.RestErr)
	DelUserUUID(tokenUUID string, accessTokenUUID string) (int64, *restDomainErr.RestErr)
}

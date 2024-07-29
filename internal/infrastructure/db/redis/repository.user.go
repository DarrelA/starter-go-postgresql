// coverage:ignore file
// Testing with integration test
package redis

import (
	"context"
	"time"

	errConst "github.com/DarrelA/starter-go-postgresql/internal/domain/error"
	restDomainErr "github.com/DarrelA/starter-go-postgresql/internal/domain/error/transport/http"
	r "github.com/DarrelA/starter-go-postgresql/internal/domain/repository/redis"
	restInterfaceErr "github.com/DarrelA/starter-go-postgresql/internal/interface/transport/http/error"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type RedisUserRepository struct {
	RedisDB *RedisDB
}

func NewUserRepository(redisDB *RedisDB) r.RedisUserRepository {
	return &RedisUserRepository{redisDB}
}

func (r RedisUserRepository) SetUserUUID(tokenUUID string, userUUID string, expiresIn int64) *restDomainErr.RestErr {
	ctx, cancel := context.WithTimeout(r.RedisDB.RedisCtx, 5*time.Second)
	defer cancel()

	timeNow := time.Now()
	err := r.RedisDB.RedisClient.Set(ctx, tokenUUID, userUUID, time.Unix(expiresIn, 0).Sub(timeNow)).Err()
	log.Error().Err(err).Msg(errConst.ErrMsgRedisError)
	return restInterfaceErr.NewInternalServerError(errConst.ErrMsgSomethingWentWrong)
}

func (r RedisUserRepository) GetUserUUID(tokenUUID string) (string, *restDomainErr.RestErr) {
	ctx, cancel := context.WithTimeout(r.RedisDB.RedisCtx, 3*time.Second)
	defer cancel()

	result, err := r.RedisDB.RedisClient.Get(ctx, tokenUUID).Result()

	if err == redis.Nil {
		err := restInterfaceErr.NewUnauthorizedError(errConst.ErrMsgPleaseLoginAgain)
		return "", err
	} else if err != nil {
		log.Error().Err(err).Msg(errConst.ErrMsgRedisError)
		return "", restInterfaceErr.NewInternalServerError(errConst.ErrMsgSomethingWentWrong)
	}

	return result, nil
}

func (r RedisUserRepository) DelUserUUID(tokenUUID string, accessTokenUUID string) (int64, *restDomainErr.RestErr) {
	ctx, cancel := context.WithTimeout(r.RedisDB.RedisCtx, 3*time.Second)
	defer cancel()
	result, err := r.RedisDB.RedisClient.Del(ctx, tokenUUID, accessTokenUUID).Result()
	if err != nil {
		log.Error().Err(err).Msg("error deleting userUUID from Redis")
		return 0, nil
	}

	return result, nil
}

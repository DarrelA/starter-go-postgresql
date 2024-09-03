// coverage:ignore file
// Testing with integration test
package redis

import (
	"context"
	"time"

	r "github.com/DarrelA/starter-go-postgresql/internal/domain/repository/redis"
	restErr "github.com/DarrelA/starter-go-postgresql/internal/error"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type RedisUserRepository struct {
	RedisDB *RedisDB
}

func NewUserRepository(redisDB *RedisDB) r.RedisUserRepository {
	return &RedisUserRepository{redisDB}
}

func (r RedisUserRepository) SetUserUUID(tokenUUID string, userUUID string, expiresIn int64) *restErr.RestErr {
	ctx, cancel := context.WithTimeout(r.RedisDB.RedisCtx, 5*time.Second)
	defer cancel()

	timeNow := time.Now()
	err := r.RedisDB.RedisClient.Set(ctx, tokenUUID, userUUID, time.Unix(expiresIn, 0).Sub(timeNow)).Err()
	if err != nil {
		log.Error().Err(err).Msg(restErr.ErrMsgRedisError)
		return restErr.NewInternalServerError(restErr.ErrMsgSomethingWentWrong)
	}

	return nil
}

func (r RedisUserRepository) GetUserUUID(tokenUUID string) (string, *restErr.RestErr) {
	ctx, cancel := context.WithTimeout(r.RedisDB.RedisCtx, 3*time.Second)
	defer cancel()

	result, err := r.RedisDB.RedisClient.Get(ctx, tokenUUID).Result()

	if err == redis.Nil {
		err := restErr.NewUnauthorizedError(restErr.ErrMsgPleaseLoginAgain)
		return "", err
	} else if err != nil {
		log.Error().Err(err).Msg(restErr.ErrMsgRedisError)
		return "", restErr.NewInternalServerError(restErr.ErrMsgSomethingWentWrong)
	}

	return result, nil
}

func (r RedisUserRepository) DelUserUUID(tokenUUID string, accessTokenUUID string) (int64, *restErr.RestErr) {
	ctx, cancel := context.WithTimeout(r.RedisDB.RedisCtx, 3*time.Second)
	defer cancel()
	result, err := r.RedisDB.RedisClient.Del(ctx, tokenUUID, accessTokenUUID).Result()
	if err != nil {
		log.Error().Err(err).Msg("error deleting userUUID from Redis")
		return 0, nil
	}

	return result, nil
}

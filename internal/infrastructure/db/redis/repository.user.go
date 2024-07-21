package redis

import (
	"context"
	"time"

	r "github.com/DarrelA/starter-go-postgresql/internal/domain/repository/redis"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

const errMsgGetUserUUID = "error_GetUserUUID"

type RedisUserRepository struct {
	RedisDB *RedisDB
}

func NewUserRepository(redisDB *RedisDB) r.RedisUserRepository {
	return &RedisUserRepository{redisDB}
}

func (r RedisUserRepository) SetUserUUID(tokenUUID string, userUUID string, expiresIn int64) error {
	ctx, cancel := context.WithTimeout(r.RedisDB.RedisCtx, 5*time.Second)
	defer cancel()

	timeNow := time.Now()
	err := r.RedisDB.RedisClient.Set(ctx, tokenUUID, userUUID, time.Unix(expiresIn, 0).Sub(timeNow)).Err()
	return err
}

func (r RedisUserRepository) GetUserUUID(tokenUUID string) (string, error) {
	ctx, cancel := context.WithTimeout(r.RedisDB.RedisCtx, 3*time.Second)
	defer cancel()

	result, err := r.RedisDB.RedisClient.Get(ctx, tokenUUID).Result()

	if err == redis.Nil {
		return "", err
	} else if err != nil {
		log.Error().Err(err).Msg(errMsgGetUserUUID)
		return "", err
	}

	return result, nil
}

func (r RedisUserRepository) DelUserUUID(tokenUUID string, accessTokenUUID string) (int64, error) {
	ctx, cancel := context.WithTimeout(r.RedisDB.RedisCtx, 3*time.Second)
	defer cancel()

	return r.RedisDB.RedisClient.Del(ctx, tokenUUID, accessTokenUUID).Result()
}

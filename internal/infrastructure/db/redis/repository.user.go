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
	redisClient *redis.Client
}

func NewUserRepository(redisClient *redis.Client) r.RedisUserRepository {
	return &RedisUserRepository{redisClient}
}

func (r RedisUserRepository) SetUserUUID(tokenUUID string, userUUID string, expiresIn int64) error {
	ctx := context.TODO()
	timeNow := time.Now()
	err := r.redisClient.Set(ctx, tokenUUID, userUUID, time.Unix(expiresIn, 0).Sub(timeNow)).Err()
	return err
}

func (r RedisUserRepository) GetUserUUID(tokenUUID string) (string, error) {
	ctx := context.TODO()
	result, err := r.redisClient.Get(ctx, tokenUUID).Result()

	if err == redis.Nil {
		return "", err
	} else if err != nil {
		log.Error().Err(err).Msg(errMsgGetUserUUID)
		return "", err
	}

	return result, nil
}

func (r RedisUserRepository) DelUserUUID(tokenUUID string, accessTokenUUID string) (int64, error) {
	ctx := context.TODO()
	return r.redisClient.Del(ctx, tokenUUID, accessTokenUUID).Result()
}

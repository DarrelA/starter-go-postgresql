package redis

import (
	"context"
	"time"

	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	repository "github.com/DarrelA/starter-go-postgresql/internal/domain/repository/redis"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type RedisDB struct {
	RedisDBConfig *entity.RedisDBConfig
	c             *redis.Client
}

func Connect(RedisDBConfig *entity.RedisDBConfig) (repository.UserRepository, error) {
	// @TODO: Switch to `context.Background()`?
	ctx := context.TODO()
	redisClient := redis.NewClient(&redis.Options{Addr: RedisDBConfig.RedisUri})
	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		log.Error().Err(err).Msg("error connecting to the Redis database")
		panic(err)
	}

	log.Info().Msg("successfully connected to the Redis database")
	return &RedisDB{RedisDBConfig: RedisDBConfig, c: redisClient}, nil
}

func (r RedisDB) SetUserUUID(tokenUUID string, userUUID string, expiresIn int64) error {
	ctx := context.TODO()
	timeNow := time.Now()
	err := r.c.Set(ctx, tokenUUID, userUUID, time.Unix(expiresIn, 0).Sub(timeNow)).Err()
	return err
}

func (r RedisDB) GetUserUUID(tokenUUID string) (string, error) {
	ctx := context.TODO()
	result, err := r.c.Get(ctx, tokenUUID).Result()

	if err == redis.Nil {
		return "", err
	} else if err != nil {
		log.Error().Err(err).Msg("error_GetUserUUID")
		return "", err
	}

	return result, nil
}

func (r RedisDB) DelUserUUID(tokenUUID string, accessTokenUUID string) (int64, error) {
	ctx := context.TODO()
	return r.c.Del(ctx, tokenUUID, accessTokenUUID).Result()
}

// @TODO: Fix Disconnect()
// func (db *RedisDB) Disconnect() {
// 	if RedisClient != nil {
// 		err := RedisClient.Close()
// 		if err != nil {
// 			log.Error().Err(err).Msg("error closing Redis database")
// 		} else {
// 			log.Info().Msg("Redis database connection closed")
// 		}
// 	}
// }

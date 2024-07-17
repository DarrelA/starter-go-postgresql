package redis

import (
	"context"

	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type RedisDB struct {
	RedisDBConfig *entity.RedisDBConfig
	RedisClient   *redis.Client
}

func Connect(redisDBConfig *entity.RedisDBConfig) (*RedisDB, error) {
	// @TODO: Switch to `context.Background()`?
	ctx := context.TODO()
	redisClient := redis.NewClient(&redis.Options{Addr: redisDBConfig.RedisUri})
	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		log.Error().Err(err).Msg("error connecting to the Redis database")
		panic(err)
	}

	log.Info().Msg("successfully connected to the Redis database")
	return &RedisDB{RedisDBConfig: redisDBConfig, RedisClient: redisClient}, nil
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

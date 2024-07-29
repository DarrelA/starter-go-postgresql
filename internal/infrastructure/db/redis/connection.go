package redis

import (
	"context"
	"time"

	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/repository"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

const (
	errMsgConnectingToDB      = "error connecting to the Redis database"
	errMsgDisconnectingFromDB = "error closing Redis database"
)

type RedisDB struct {
	RedisDBConfig *entity.RedisDBConfig
	RedisClient   *redis.Client
	RedisCtx      context.Context
}

func (r *RedisDB) ConnectToRedis(redisDBConfig *entity.RedisDBConfig) repository.InMemoryDB {
	// Create a top level context
	ctx := context.Background()

	redisClient := redis.NewClient(&redis.Options{Addr: redisDBConfig.RedisUri})
	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		log.Error().Err(err).Msg(errMsgConnectingToDB)
		panic(err)
	}

	log.Info().Msg("successfully connected to the Redis database")
	return &RedisDB{RedisDBConfig: redisDBConfig, RedisClient: redisClient, RedisCtx: ctx}
}

func (r *RedisDB) DisconnectFromRedis() {
	if r != nil {
		_, cancel := context.WithTimeout(r.RedisCtx, 10*time.Second)
		defer cancel()

		err := r.RedisClient.Close()
		if err != nil {
			log.Error().Err(err).Msg(errMsgDisconnectingFromDB)
		} else {
			log.Info().Msg("Redis database connection closed")
		}
	}
}

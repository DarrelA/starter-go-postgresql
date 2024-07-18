package redis

import (
	"context"

	"github.com/DarrelA/starter-go-postgresql/internal/application/repository"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type RedisDB struct {
	RedisDBConfig *entity.RedisDBConfig
	RedisClient   *redis.Client
}

// Connection is a struct to hold the return values from the `Connect` function.
type Connection struct {
	InMemoryDB repository.InMemoryDB
	RedisDB    *RedisDB
}

func Connect(redisDBConfig *entity.RedisDBConfig) Connection {
	// @TODO: Switch to `context.Background()`?
	ctx := context.TODO()
	redisClient := redis.NewClient(&redis.Options{Addr: redisDBConfig.RedisUri})
	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		log.Error().Err(err).Msg("error connecting to the Redis database")
		panic(err)
	}

	log.Info().Msg("successfully connected to the Redis database")
	redisDB := &RedisDB{RedisDBConfig: redisDBConfig, RedisClient: redisClient}
	return Connection{InMemoryDB: redisDB, RedisDB: redisDB}
}

func (r *RedisDB) Disconnect() {
	if r != nil {
		err := r.RedisClient.Close()
		if err != nil {
			log.Error().Err(err).Msg("error closing Redis database")
		} else {
			log.Info().Msg("Redis database connection closed")
		}
	}
}

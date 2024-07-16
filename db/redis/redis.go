package db

import (
	"context"

	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type RedisDB struct {
	config *entity.RedisDBConfig
}

var (
	RedisClient *redis.Client
	ctx         context.Context
)

// NewDB creates a new RedisDB instance with loaded config
func NewDB(config *entity.RedisDBConfig) *RedisDB {
	return &RedisDB{
		config: config,
	}
}

func (db *RedisDB) Connect(config *entity.RedisDBConfig) {
	ctx = context.TODO()

	RedisClient = redis.NewClient(&redis.Options{
		Addr: config.RedisUri,
	})

	if _, err := RedisClient.Ping(ctx).Result(); err != nil {
		log.Panic().Err(err).Msg("error connecting to the Redis database")
		panic(err)
	}

	log.Info().Msg("successfully connected to the Redis database")
}

func (db *RedisDB) Disconnect() {
	if RedisClient != nil {
		err := RedisClient.Close()
		if err != nil {
			log.Error().Err(err).Msg("error closing Redis database")
		} else {
			log.Info().Msg("Redis database connection closed")
		}
	}
}

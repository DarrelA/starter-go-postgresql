package db

import (
	"context"

	"github.com/DarrelA/starter-go-postgresql/configs"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type RedisDB struct {
	configs.RedisDBConfig
}

var (
	RedisClient *redis.Client
	ctx         context.Context
)

// NewDB creates a new RedisDB instance with loaded config
func NewDB() *RedisDB {
	return &RedisDB{
		RedisDBConfig: configs.RedisDB,
	}
}

func (db *RedisDB) Connect() {
	dbCfg := configs.RedisDB
	ctx = context.TODO()

	RedisClient = redis.NewClient(&redis.Options{
		Addr: dbCfg.RedisUri,
	})

	if _, err := RedisClient.Ping(ctx).Result(); err != nil {
		panic(err)
	}

	log.Info().Msg("successfully connected to the Redis database")
}

func (db *RedisDB) Disconnect() {
	if RedisClient != nil {
		err := RedisClient.Close()
		if err != nil {
			log.Error().Msg("error closing Redis database: " + err.Error())
		} else {
			log.Info().Msg("Redis database connection closed")
		}
	}
}

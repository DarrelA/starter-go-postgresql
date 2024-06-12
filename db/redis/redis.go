package db

import (
	"context"
	"log"

	"github.com/DarrelA/starter-go-postgresql/configs"
	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
	ctx         context.Context
)

func ConnectRedis() {
	dbCfg := configs.RedisDB
	ctx = context.TODO()

	RedisClient = redis.NewClient(&redis.Options{
		Addr: dbCfg.RedisUri,
	})

	if _, err := RedisClient.Ping(ctx).Result(); err != nil {
		panic(err)
	}

	log.Println("Successfully connected to the Redis database")
}

package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	"github.com/redis/go-redis/v9"
)

const errMsgConnectingToDB = "error connecting to the Redis database"

// Open creates and verifies a Redis client. The caller owns the returned client
// and must close it.
func Open(ctx context.Context, redisDBConfig *entity.RedisDBConfig) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{Addr: redisDBConfig.RedisUri})
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if _, err := redisClient.Ping(pingCtx).Result(); err != nil {
		_ = redisClient.Close()
		return nil, fmt.Errorf("%s: %w", errMsgConnectingToDB, err)
	}

	return redisClient, nil
}

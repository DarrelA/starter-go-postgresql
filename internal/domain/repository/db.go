package repository

import "github.com/DarrelA/starter-go-postgresql/internal/domain/entity"

type RDBMS interface {
	ConnectToPostgres(postgresDBConfig *entity.PostgresDBConfig) RDBMS
	DisconnectFromPostgres()
}

type InMemoryDB interface {
	ConnectToRedis(redisDBConfig *entity.RedisDBConfig) InMemoryDB
	DisconnectFromRedis()
}

package db

import "github.com/DarrelA/starter-go-postgresql/internal/domain/entity"

type RDBMS interface {
	Connect()
	Disconnect()
}

type InMemoryDB interface {
	Connect(config *entity.RedisDBConfig)
	Disconnect()
}

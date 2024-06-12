package app

import (
	"github.com/DarrelA/starter-go-postgresql/configs"
	pgdb "github.com/DarrelA/starter-go-postgresql/db/pgdb"
	redis "github.com/DarrelA/starter-go-postgresql/db/redis"
	"github.com/gofiber/fiber/v2"
)

var (
	app = fiber.New()
)

func StartApp() {
	pgdb.ConnectPostgresDatabase()
	redis.ConnectRedis()
	mapUrls()
	app.Listen(":" + configs.Port)
}

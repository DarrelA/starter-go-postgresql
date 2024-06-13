package app

import (
	"log"

	"github.com/DarrelA/starter-go-postgresql/configs"
	pgdb "github.com/DarrelA/starter-go-postgresql/db/pgdb"
	redisDb "github.com/DarrelA/starter-go-postgresql/db/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

var (
	app         = fiber.New()
	authService = fiber.New()
)

func StartApp() {
	app.Mount("/api", authService)
	app.Use(logger.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins:     configs.CORSSettings.AllowedOrigins,
		AllowMethods:     "GET,POST",
		AllowHeaders:     "Content-Type",
		ExposeHeaders:    "Content-Length",
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60,
	}))

	pgdb.ConnectPostgresDatabase()
	redisDb.ConnectRedis()
	authServiceRouter()
	log.Fatal(app.Listen(":" + configs.Port))
}

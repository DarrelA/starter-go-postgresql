package app

import (
	"github.com/DarrelA/starter-go-postgresql/configs"
	"github.com/DarrelA/starter-go-postgresql/db"
	redisDb "github.com/DarrelA/starter-go-postgresql/db/redis"
	"github.com/DarrelA/starter-go-postgresql/internal/middlewares"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rs/zerolog/log"
)

var (
	app         = fiber.New()
	authService = fiber.New()
)

func StartApp(rdbms db.RDBMS) {
	pgdb := rdbms

	app.Mount("/auth", authService)

	// Middlewares
	authService.Use(cors.New(cors.Config{
		AllowOrigins:     configs.CORSSettings.AllowedOrigins,
		AllowMethods:     "GET,POST",
		AllowHeaders:     "Content-Type",
		ExposeHeaders:    "Content-Length",
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60,
	}))
	authService.Use(middlewares.CorrelationAndRequestID)
	authService.Use(middlewares.LogRequest)

	pgdb.Connect()
	redisDb.ConnectRedis()
	authServiceRouter(authService)

	err := app.Listen(":" + configs.Port)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}
}

func CloseConnections(rdbms db.RDBMS) {
	pgdb := rdbms

	pgdb.Disconnect()
	redisDb.DisconnectRedis()
}

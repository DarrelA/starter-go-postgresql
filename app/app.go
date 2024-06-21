package app

import (
	"context"
	"time"

	"github.com/DarrelA/starter-go-postgresql/configs"
	"github.com/DarrelA/starter-go-postgresql/db"
	"github.com/DarrelA/starter-go-postgresql/internal/middlewares"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rs/zerolog/log"
)

func StartApp(rdbms db.RDBMS, inMemoryDb db.InMemoryDB) *fiber.App {
	log.Debug().Msg("starting StartApp function")
	app := fiber.New()
	authService := fiber.New()
	log.Debug().Msg("created app and authService instance")

	pgdb := rdbms
	redisDb := inMemoryDb
	pgdb.Connect()
	redisDb.Connect()
	log.Debug().Msg("connected to databases")

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
	log.Debug().Msg("applied middlewares to authService")

	authServiceRouter(authService)

	return app
}

func CloseConnections(app *fiber.App, rdbms db.RDBMS, inMemoryDb db.InMemoryDB) {
	log.Debug().Msg("starting CloseConnections function")

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Err(err).Msg("failed to gracefully shutdown the server")
	}
	cancel()
	log.Debug().Msg("app instance has shutdown")

	pgdb := rdbms
	redisDb := inMemoryDb

	pgdb.Disconnect()
	redisDb.Disconnect()
	log.Debug().Msg("end of CloseConnections function")
}

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
	createDBConnections(rdbms, inMemoryDb)
	app, authService := createFiberInstances()
	useMiddleware(authService)
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

func createDBConnections(rdbms db.RDBMS, inMemoryDb db.InMemoryDB) {
	pgdb := rdbms
	redisDb := inMemoryDb
	pgdb.Connect()
	redisDb.Connect()
	log.Debug().Msg("connected to databases")
}

func createFiberInstances() (*fiber.App, *fiber.App) {
	app := fiber.New()
	authService := fiber.New()
	log.Debug().Msg("created app and authService instance")
	app.Mount("/auth", authService)
	return app, authService
}

func useMiddleware(authService *fiber.App) {
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
}

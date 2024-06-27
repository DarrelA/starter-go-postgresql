package app

import (
	"context"
	"time"

	"github.com/DarrelA/starter-go-postgresql/configs"
	"github.com/DarrelA/starter-go-postgresql/db"
	pgdb "github.com/DarrelA/starter-go-postgresql/db/pgdb"
	redisDb "github.com/DarrelA/starter-go-postgresql/db/redis"
	"github.com/DarrelA/starter-go-postgresql/internal/middlewares"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rs/zerolog/log"
)

var (
	rdbmsInstance       db.RDBMS
	inMemoryDbInstance  db.InMemoryDB
	appInstance         *fiber.App
	authServiceInstance *fiber.App
)

func CreateDBConnections() (db.RDBMS, db.InMemoryDB) {
	if rdbmsInstance == nil {
		rdbmsInstance = pgdb.NewDB()
	}

	if inMemoryDbInstance == nil {
		inMemoryDbInstance = redisDb.NewDB()
	}

	rdbmsInstance.Connect()
	inMemoryDbInstance.Connect()

	return rdbmsInstance, inMemoryDbInstance
}

func CloseConnections() {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	if err := appInstance.ShutdownWithContext(ctx); err != nil {
		log.Err(err).Msg("failed to gracefully shutdown the server")
	}

	cancel()
	log.Info().Msg("app instance has shutdown")

	rdbmsInstance.Disconnect()
	inMemoryDbInstance.Disconnect()
}

// ConfigureAppInstance sets up and configures instances of Fiber for the main app and auth service,
// including middleware and routing for authentication.
func ConfigureAppInstance() (*fiber.App, *fiber.App) {
	log.Info().Msg("creating fiber instances, connecting middleware & router")
	appInstance = fiber.New()
	authServiceInstance = fiber.New()
	appInstance.Mount("/auth", authServiceInstance)

	useMiddlewares(authServiceInstance)
	authServiceRouter(authServiceInstance)

	log.Debug().Msgf("appInstance memory address: %p", appInstance)
	return appInstance, authServiceInstance
}

func useMiddlewares(authServiceInstance *fiber.App) {
	authServiceInstance.Use(cors.New(cors.Config{
		AllowOrigins:     configs.CORSSettings.AllowedOrigins,
		AllowMethods:     "GET,POST",
		AllowHeaders:     "Content-Type",
		ExposeHeaders:    "Content-Length",
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60,
	}))

	authServiceInstance.Use(middlewares.CorrelationAndRequestID)
	authServiceInstance.Use(middlewares.LogRequest)
	log.Info().Msg("applied middlewares to authServiceInstance")
}

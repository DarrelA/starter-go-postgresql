package app

import (
	"context"
	"runtime/debug"
	"time"

	"github.com/DarrelA/starter-go-postgresql/configs"
	"github.com/DarrelA/starter-go-postgresql/db"
	pgdb "github.com/DarrelA/starter-go-postgresql/db/pgdb"
	redisDb "github.com/DarrelA/starter-go-postgresql/db/redis"
	"github.com/DarrelA/starter-go-postgresql/internal/middlewares"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
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

func StartServer() {
	log.Info().Msg("listening at port: " + configs.Port)
	err := appInstance.Listen(":" + configs.Port)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}
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

func useMiddlewares(authServiceInstance *fiber.App) {
	// Recover middleware to catch panics and handle errors
	authServiceInstance.Use(recover.New(recover.Config{
		EnableStackTrace:  true,
		StackTraceHandler: customStackTraceHandler,
	}))

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

func customStackTraceHandler(c *fiber.Ctx, e interface{}) {
	stackTrace := string(debug.Stack())

	// Log the panic and stack trace
	if err, ok := e.(error); ok {
		log.Error().
			Err(err).
			Str("stack_trace", stackTrace).
			Msg("middleware_error")
	} else {
		// The Interface method is used to log the panic value itself, which could be of any type.
		log.Error().
			Interface("error", e).
			Str("stack_trace", stackTrace).
			Msg("middleware_error")
	}

	c.Status(fiber.StatusServiceUnavailable).
		JSON(fiber.Map{"status": "fail", "message": "service is unavailable at the moment"})
}

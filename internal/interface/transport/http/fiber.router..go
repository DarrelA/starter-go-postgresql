package http

import (
	"runtime/debug"

	"github.com/DarrelA/starter-go-postgresql/configs"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/factory"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/service"
	mw "github.com/DarrelA/starter-go-postgresql/internal/interface/middleware"
	"github.com/DarrelA/starter-go-postgresql/internal/utils/err_rest"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog/log"
)

func StartServer(app *fiber.App) {
	log.Info().Msg("listening at port: " + configs.Port)
	err := app.Listen(":" + configs.Port)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}
}

func NewRouter(
	token service.TokenService,
	userFactory factory.UserFactory,
	authHandler AuthHandler,
) *fiber.App {
	log.Info().Msg("creating fiber instances")
	appInstance := fiber.New()
	authServiceInstance := fiber.New()
	appInstance.Mount("/auth", authServiceInstance)

	log.Info().Msg("connecting middlewares")
	useMiddlewares(authServiceInstance)

	log.Info().Msg("setting up routes")
	v1 := authServiceInstance.Group("/api/v1", func(c *fiber.Ctx) error { // middleware for /api/v1
		c.Set("Version", "v1")
		return c.Next()
	})

	user := v1.Group("/users")
	user.Post("/register", mw.PreProcessInputs, authHandler.Register)
	user.Post("/login", mw.PreProcessInputs, authHandler.Login)

	authUser := user.Group("/").Use(mw.Deserializer(token, userFactory))
	authUser.Get("/logout", authHandler.Logout)

	user.Get("/refresh", authHandler.RefreshAccessToken)
	// @TODO: Fix route
	// user.Get("/users/me", mw.Deserializer, handlers.GetMe)

	authServiceInstance.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
	})

	appInstance.All("*", func(c *fiber.Ctx) error {
		path := c.Path()
		err := err_rest.NewBadRequestError("Invalid Path: " + path)
		log.Error().Err(err).Msg("")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "404 - Not Found",
		})
	})

	log.Info().Msg("/health endpoint is available")
	log.Debug().Msgf("appInstance memory address: %p", appInstance)
	log.Debug().Msgf("authServiceInstance memory address: %p", authServiceInstance)
	return appInstance
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

	authServiceInstance.Use(mw.CorrelationAndRequestID)
	authServiceInstance.Use(mw.LoggerMW)
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

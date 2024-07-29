package http

import (
	"runtime/debug"

	"github.com/DarrelA/starter-go-postgresql/internal/application/factory"
	appSvc "github.com/DarrelA/starter-go-postgresql/internal/application/service"
	r "github.com/DarrelA/starter-go-postgresql/internal/domain/repository/redis"
	domainSvc "github.com/DarrelA/starter-go-postgresql/internal/domain/service"
	"github.com/DarrelA/starter-go-postgresql/internal/infrastructure/config"
	mw "github.com/DarrelA/starter-go-postgresql/internal/interface/middleware"
	ppmw "github.com/DarrelA/starter-go-postgresql/internal/interface/middleware/preprocess_inputs"
	restInterfaceErr "github.com/DarrelA/starter-go-postgresql/internal/interface/transport/http/error"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog/log"
)

const (
	errMsgStartServerFailure = "failed to start server"
	errMsgServiceUnavailable = "service is unavailable at the moment"
)

func StartServer(app *fiber.App, port string) {
	log.Info().Msg("listening at port: " + port)
	err := app.Listen(":" + port)
	if err != nil {
		log.Error().Err(err).Msg(errMsgStartServerFailure)
	}
}

func NewRouter(
	envConfig *config.EnvConfig,
	redisRepo r.RedisUserRepository,
	tokenService domainSvc.TokenService,
	userFactory factory.UserFactory,
	authService appSvc.AuthService,
	userService appSvc.UserService,
) *fiber.App {
	log.Info().Msg("creating fiber instances")
	appInstance := fiber.New()
	authServiceInstance := fiber.New()
	appInstance.Mount("/auth", authServiceInstance)

	log.Info().Msg("connecting middlewares")
	useMiddlewares(authServiceInstance, envConfig)

	log.Info().Msg("setting up routes")
	v1 := authServiceInstance.Group("/api/v1", func(c *fiber.Ctx) error { // middleware for /api/v1
		c.Set("Version", "v1")
		return c.Next()
	})

	user := v1.Group("/users")
	user.Post("/register", ppmw.PreProcessInputs, authService.Register)
	user.Post("/login", ppmw.PreProcessInputs, authService.Login)

	authUser := user.Group("/").Use(mw.Deserializer(redisRepo, tokenService, userFactory))
	authUser.Get("/logout", authService.Logout)
	authUser.Get("/me", userService.GetUserRecord)

	user.Get("/refresh", authService.RefreshAccessToken)

	authServiceInstance.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
	})

	appInstance.All("*", func(c *fiber.Ctx) error {
		path := c.Path()
		err := restInterfaceErr.NewBadRequestError("Invalid Path: " + path)
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

func useMiddlewares(authServiceInstance *fiber.App, envConfig *config.EnvConfig) {
	// Recover middleware to catch panics and handle errors
	authServiceInstance.Use(recover.New(recover.Config{
		EnableStackTrace:  true,
		StackTraceHandler: customStackTraceHandler,
	}))

	authServiceInstance.Use(func(c *fiber.Ctx) error {
		c.Locals("baseURLsConfig", envConfig.BaseURLsConfig)
		c.Locals("env", envConfig.Env)
		return c.Next()
	})

	authServiceInstance.Use(cors.New(cors.Config{
		AllowOrigins:     envConfig.CORSConfig.AllowedOrigins,
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
		JSON(fiber.Map{"status": "fail", "message": errMsgServiceUnavailable})
}

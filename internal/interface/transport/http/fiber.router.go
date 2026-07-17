package http

import (
	"runtime/debug"

	"github.com/DarrelA/starter-go-postgresql/internal/domain/repository"
	domainSvc "github.com/DarrelA/starter-go-postgresql/internal/domain/service"
	restErr "github.com/DarrelA/starter-go-postgresql/internal/error"
	"github.com/DarrelA/starter-go-postgresql/internal/infrastructure/config"
	mw "github.com/DarrelA/starter-go-postgresql/internal/interface/middleware"
	authmw "github.com/DarrelA/starter-go-postgresql/internal/interface/middleware/deserialize_user"
	ppmw "github.com/DarrelA/starter-go-postgresql/internal/interface/middleware/preprocess_inputs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog/log"
)

const (
	errMsgStartServerFailure = "failed to start server"
	errMsgServiceUnavailable = "service is unavailable at the moment"
)

// StartServer listens on port and blocks until the Fiber server stops.
func StartServer(app *fiber.App, port string) error {
	log.Info().Msg("listening at port: " + port)
	return app.Listen(":" + port)
}

// AuthRoutes defines the authentication handlers required by NewRouter.
type AuthRoutes interface {
	Register(*fiber.Ctx) error
	Login(*fiber.Ctx) error
	RefreshAccessToken(*fiber.Ctx) error
	Logout(*fiber.Ctx) error
}

// UserRoutes defines the current-user handlers required by NewRouter.
type UserRoutes interface {
	GetUserRecord(*fiber.Ctx) error
}

// OAuthRoutes defines the OAuth handlers required by NewRouter.
type OAuthRoutes interface {
	Login(*fiber.Ctx) error
	Callback(*fiber.Ctx) error
}

// NewRouter constructs the Fiber application and connects its middleware and routes.
func NewRouter(
	envConfig *config.EnvConfig,
	redisRepo repository.TokenRepository,
	tokenService domainSvc.TokenService,
	userUseCase UserRoutes,
	authUseCase AuthRoutes,
	googleOAuth2UseCase OAuthRoutes,
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

	/********************
	 *   Refresh Token  *
	 ********************/
	user := v1.Group("/users")
	user.Post("/register", ppmw.PreProcessInputs, authUseCase.Register)
	user.Post("/login", ppmw.PreProcessInputs, authUseCase.Login)

	user.Post("/logout", authUseCase.Logout)
	user.Post("/refresh", authUseCase.RefreshAccessToken)
	user.Get("/me", authmw.Authenticate(redisRepo, tokenService), userUseCase.GetUserRecord)

	/********************
	 *      OAuth2      *
	 ********************/
	authServiceInstance.Get("/google_login", googleOAuth2UseCase.Login)
	authServiceInstance.Get("/google_callback", googleOAuth2UseCase.Callback)

	authServiceInstance.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
	})

	appInstance.All("*", func(c *fiber.Ctx) error {
		path := c.Path()
		err := restErr.NewBadRequestError("Invalid Path: " + path)
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
		AllowHeaders:     "Content-Type, Authorization, X-CSRF-Token",
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

package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/DarrelA/starter-go-postgresql/internal/application/factory"
	"github.com/DarrelA/starter-go-postgresql/internal/infrastructure/config"
	"github.com/DarrelA/starter-go-postgresql/internal/infrastructure/db/postgres"
	"github.com/DarrelA/starter-go-postgresql/internal/infrastructure/db/redis"
	jwt "github.com/DarrelA/starter-go-postgresql/internal/infrastructure/jwt/service"
	logger "github.com/DarrelA/starter-go-postgresql/internal/infrastructure/logger/zerolog"
	"github.com/DarrelA/starter-go-postgresql/internal/interface/transport/http"
	"github.com/rs/zerolog/log"
)

func main() {
	logger.CreateAppLog()
	startApp()
	waitForShutdown()
	logger.GetLogFile().Close()
	os.Exit(0)
}

func startApp() {
	// Use `WaitGroup` when you just need to wait for tasks to complete without exchanging data.
	// Use channels when you need to signal task completion and possibly exchange data.
	var wg sync.WaitGroup
	wg.Add(1)

	// Functions are running synchronously
	go func() {
		defer wg.Done()

		envConfig := config.NewTokenService()
		envConfig.LoadAppConfig()
		envConfig.LoadLogConfig()
		envConfig.LoadDBConfig()
		envConfig.LoadRedisConfig()
		envConfig.LoadJWTConfig()
		envConfig.LoadCORSConfig()

		if c, ok := envConfig.(*config.EnvConfig); ok {
			redisDB, err := redis.Connect(c.RedisDBConfig)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to connect to redis")
			}

			redisUserRepo := redis.NewUserRepository(redisDB.RedisClient)

			postgresDB, err := postgres.Connect(c.PostgresDBConfig)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to connect to postgres")
			}

			postgresUserRepo := postgres.NewUserRepository(postgresDB.Dbpool)

			// Dependency injection
			// User
			userFactory := factory.NewUserFactory(c.JWTConfig, postgresUserRepo)
			userService := http.NewUserService()

			// Token
			tokenService := jwt.NewTokenService()

			// Auth
			authService := http.NewAuthService(redisUserRepo, userFactory, tokenService)

			// app.SeedDatabase(postgresUserRepo)
			appServiceInstance := http.NewRouter(
				c, redisUserRepo,
				tokenService, userFactory, authService, userService,
			)
			go http.StartServer(appServiceInstance, c.Port)

		} else {
			log.Fatal().Msg("failed to load environment configuration")
		}
	}()

	wg.Wait()
}

func waitForShutdown() {
	sigChan := make(chan os.Signal, 1) // Create a channel to listen for OS signals
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan // Block until a signal is received
	log.Debug().Msg("received termination signal, shutting down")
	//TODO: Decouple Redis
	// app.CloseConnections()
}

package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	rp "github.com/DarrelA/starter-go-postgresql/internal/domain/repository/postgres"
	rr "github.com/DarrelA/starter-go-postgresql/internal/domain/repository/redis"
	"github.com/DarrelA/starter-go-postgresql/internal/infrastructure/config"
	"github.com/DarrelA/starter-go-postgresql/internal/infrastructure/db/postgres"
	"github.com/DarrelA/starter-go-postgresql/internal/infrastructure/db/redis"
	jwt "github.com/DarrelA/starter-go-postgresql/internal/infrastructure/jwt/service"
	envLogger "github.com/DarrelA/starter-go-postgresql/internal/infrastructure/logger"
	logger "github.com/DarrelA/starter-go-postgresql/internal/infrastructure/logger/zerolog"
	"github.com/DarrelA/starter-go-postgresql/internal/interface/factory"
	"github.com/DarrelA/starter-go-postgresql/internal/interface/transport/http"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func main() {
	logFile := envLogger.CreateAppLog()
	logger.NewZeroLogger(logFile)
	config := initializeEnv()
	rc, redisUserRepo, pc, postgresUserRepo := initializeDatabases(config)

	// Use `WaitGroup` when you just need to wait for tasks to complete without exchanging data.
	// Use channels when you need to signal task completion and possibly exchange data.
	var wg sync.WaitGroup
	wg.Add(1)
	appServiceInstance := initializeServer(&wg, config, redisUserRepo, postgresUserRepo)

	wg.Wait()
	waitForShutdown(appServiceInstance, rc, pc)
	logFile.Close()
	os.Exit(0)
}

func initializeEnv() *config.EnvConfig {
	envConfig := config.NewTokenService()
	envConfig.LoadAppConfig()
	envConfig.LoadLogConfig()
	envConfig.LoadDBConfig()
	envConfig.LoadRedisConfig()
	envConfig.LoadJWTConfig()
	envConfig.LoadCORSConfig()
	config, ok := envConfig.(*config.EnvConfig)
	if !ok {
		log.Fatal().Msg("failed to load environment configuration")
	}

	return config
}

func initializeDatabases(config *config.EnvConfig) (
	redis.Connection, rr.RedisUserRepository,
	postgres.Connection, rp.PostgresUserRepository,
) {
	redisConnection := redis.Connect(config.RedisDBConfig)
	redisUserRepo := redis.NewUserRepository(redisConnection.RedisDB)

	postgresConnection := postgres.Connect(config.PostgresDBConfig)
	postgresUserRepo := postgres.NewUserRepository(postgresConnection.PostgresDB.Dbpool)
	postgresSeedRepo := postgres.NewSeedRepository(postgresConnection.PostgresDB.Dbpool, config.Env)
	postgresSeedRepo.Seed(postgresUserRepo)

	return redisConnection, redisUserRepo, postgresConnection, postgresUserRepo
}

func initializeServer(
	wg *sync.WaitGroup, config *config.EnvConfig,
	redisUserRepo rr.RedisUserRepository, postgresUserRepo rp.PostgresUserRepository,
) *fiber.App {
	defer wg.Done()
	userFactory := factory.NewUserFactory(config.JWTConfig, postgresUserRepo)
	userService := http.NewUserService()
	tokenService := jwt.NewTokenService()
	authService := http.NewAuthService(redisUserRepo, userFactory, tokenService)

	appServiceInstance := http.NewRouter(config, redisUserRepo,
		tokenService, userFactory, authService, userService,
	)

	go func() {
		http.StartServer(appServiceInstance, config.Port)
	}()
	return appServiceInstance
}

func waitForShutdown(appServiceInstance *fiber.App, rc redis.Connection, pc postgres.Connection) {
	sigChan := make(chan os.Signal, 1) // Create a channel to listen for OS signals
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan // Block until a signal is received
	log.Debug().Msg("received termination signal, shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	if err := appServiceInstance.ShutdownWithContext(ctx); err != nil {
		log.Err(err).Msg("failed to gracefully shutdown the server")
	}

	cancel()
	log.Info().Msg("app instance has shutdown")

	rc.InMemoryDB.Disconnect()
	pc.RDBMS.Disconnect()
}

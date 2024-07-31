// coverage:ignore file
// Testing with integration test
package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	repo "github.com/DarrelA/starter-go-postgresql/internal/domain/repository"
	rp "github.com/DarrelA/starter-go-postgresql/internal/domain/repository/postgres"
	rr "github.com/DarrelA/starter-go-postgresql/internal/domain/repository/redis"

	"github.com/DarrelA/starter-go-postgresql/internal/infrastructure/config"
	"github.com/DarrelA/starter-go-postgresql/internal/infrastructure/db/postgres"
	"github.com/DarrelA/starter-go-postgresql/internal/infrastructure/db/redis"
	jwt "github.com/DarrelA/starter-go-postgresql/internal/infrastructure/jwt"
	envLogger "github.com/DarrelA/starter-go-postgresql/internal/infrastructure/logger"
	logger "github.com/DarrelA/starter-go-postgresql/internal/infrastructure/logger/zerolog"

	interfaceSvc "github.com/DarrelA/starter-go-postgresql/internal/interface/service"
	"github.com/DarrelA/starter-go-postgresql/internal/interface/transport/http"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func main() {
	envLogger.LogCWD()
	envLogger.ListFiles()
	logFile := envLogger.CreateAppLog("/docker_wd/logs/app.log")
	logger.NewZeroLogger(logFile)
	config := initializeEnv()
	redisConn, redisUserRepo, postgresConn, postgresUserRepo := initializeDatabases(config)

	// Use `WaitGroup` when you just need to wait for tasks to complete without exchanging data.
	// Use channels when you need to signal task completion and possibly exchange data.
	var wg sync.WaitGroup
	wg.Add(1)
	appServiceInstance := initializeServer(&wg, config, redisUserRepo, postgresUserRepo)

	wg.Wait()
	waitForShutdown(appServiceInstance, redisConn, postgresConn)
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
		log.Error().Msg("failed to load environment configuration")
	}

	return config
}

func initializeDatabases(config *config.EnvConfig) (
	repo.InMemoryDB, rr.RedisUserRepository,
	repo.RDBMS, rp.PostgresUserRepository,
) {
	redisDB := &redis.RedisDB{}
	redisConnection := redisDB.ConnectToRedis(config.RedisDBConfig)
	redisDBInstance := redisConnection.(*redis.RedisDB) // Type assert redisDB to *redis.RedisDB
	redisUserRepo := redis.NewUserRepository(redisDBInstance)

	postgresDB := &postgres.PostgresDB{}
	postgresConnection := postgresDB.ConnectToPostgres(config.PostgresDBConfig)
	postgresDBInstance := postgresConnection.(*postgres.PostgresDB) // Type assert postgresDB to *postgres.PostgresDB
	postgresUserRepo := postgres.NewUserRepository(postgresDBInstance.Dbpool)
	postgresSeedRepo := postgres.NewSeedRepository(postgresDBInstance.Dbpool, config.Env)
	postgresSeedRepo.Seed(postgresUserRepo)

	return redisDBInstance, redisUserRepo, postgresConnection, postgresUserRepo
}

func initializeServer(
	wg *sync.WaitGroup, config *config.EnvConfig,
	redisUserRepo rr.RedisUserRepository, postgresUserRepo rp.PostgresUserRepository,
) *fiber.App {
	defer wg.Done()
	userService := interfaceSvc.NewUserService(config.JWTConfig, postgresUserRepo)
	userUseCase := http.NewUserUseCase()
	tokenService := jwt.NewTokenService()
	authUseCase := http.NewAuthUseCase(redisUserRepo, userService, tokenService)

	appServiceInstance := http.NewRouter(config, redisUserRepo,
		tokenService, userService, authUseCase, userUseCase,
	)

	go func() {
		http.StartServer(appServiceInstance, config.Port)
	}()
	return appServiceInstance
}

func waitForShutdown(appServiceInstance *fiber.App, redisConn repo.InMemoryDB, postgresConn repo.RDBMS) {
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

	redisConn.DisconnectFromRedis()
	postgresConn.DisconnectFromPostgres()
}

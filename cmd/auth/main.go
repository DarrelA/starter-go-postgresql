package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/DarrelA/starter-go-postgresql/app"
	"github.com/DarrelA/starter-go-postgresql/configs"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/factory"
	"github.com/DarrelA/starter-go-postgresql/internal/infrastructure/db/postgres"
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

		// @TODO: Refactor the codes for Redis
		app.CreateRedisConnection()

		dbpool, err := postgres.NewRDBMS(&configs.PGDB)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to connect to the database")
		}

		// Dependency injection
		// User
		userRepo := postgres.NewUserRepository(dbpool)
		// @TODO: Ensure interface compliance
		userFactory := factory.NewUserFactory(userRepo)
		userService := http.NewUserService()

		// Token
		tokenService := jwt.NewTokenService()

		// Auth
		// @TODO: Ensure interface compliance
		authService := http.NewAuthService(*userFactory, tokenService)

		app.SeedDatabase(dbpool)
		// appServiceInstance := http.NewRouter(tokenService, *userFactory, *authService)
		appServiceInstance := http.NewRouter(tokenService, *userFactory, *authService, userService)
		go http.StartServer(appServiceInstance)
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

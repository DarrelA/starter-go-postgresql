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
	"github.com/DarrelA/starter-go-postgresql/internal/interface/transport/http"
	envs_utils "github.com/DarrelA/starter-go-postgresql/internal/utils/envs"
	"github.com/rs/zerolog/log"
)

func main() {
	envs_utils.CreateAppLog()
	startApp()
	waitForShutdown()
	envs_utils.GetLogFile().Close()
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
		userFactory := factory.NewUserFactory(userRepo)

		// Token
		tokenService := jwt.NewTokenService()

		// Auth
		authService := http.NewAuthService(*userFactory, tokenService)

		app.SeedDatabase(dbpool)
		authServiceInstance := http.NewRouter(tokenService, *userFactory, *authService)
		go http.StartServer(authServiceInstance)
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

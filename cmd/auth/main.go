package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	appservice "github.com/DarrelA/starter-go-postgresql/internal/application/service"
	"github.com/DarrelA/starter-go-postgresql/internal/infrastructure/auth"
	"github.com/DarrelA/starter-go-postgresql/internal/infrastructure/config"
	"github.com/DarrelA/starter-go-postgresql/internal/infrastructure/db/postgres"
	"github.com/DarrelA/starter-go-postgresql/internal/infrastructure/db/redis"
	envLogger "github.com/DarrelA/starter-go-postgresql/internal/infrastructure/logger"
	logger "github.com/DarrelA/starter-go-postgresql/internal/infrastructure/logger/zerolog"
	httptransport "github.com/DarrelA/starter-go-postgresql/internal/interface/transport/http"
	oauth2transport "github.com/DarrelA/starter-go-postgresql/internal/interface/transport/http/oauth2"
	"github.com/rs/zerolog/log"
)

func main() {
	if err := run(); err != nil {
		log.Fatal().Err(err).Msg("application stopped")
	}
}

func run() error {
	envConfig, err := config.LoadEnvConfig()
	if err != nil {
		return err
	}
	if err := config.ConfigureLogging(envConfig.LogLevel); err != nil {
		return err
	}
	tokenService, err := auth.NewJWTService(envConfig.JWTConfig)
	if err != nil {
		return fmt.Errorf("initialize token service: %w", err)
	}
	csrfService, err := auth.NewCSRFService(envConfig.CSRFConfig.Secret)
	if err != nil {
		return fmt.Errorf("initialize CSRF service: %w", err)
	}

	envLogger.LogCWD()
	envLogger.ListFiles()
	logFile := envLogger.CreateAppLog("/docker_wd/logs/app.log")
	if logFile != nil {
		defer logFile.Close()
		logger.NewZeroLogger(logFile)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	postgresPool, err := postgres.Open(ctx, envConfig.PostgresDBConfig)
	if err != nil {
		return err
	}
	defer postgresPool.Close()
	if envConfig.MigrateOnStartup {
		if err := postgres.Migrate(ctx, postgresPool); err != nil {
			return err
		}
	}

	redisClient, err := redis.Open(ctx, envConfig.RedisDBConfig)
	if err != nil {
		return err
	}
	defer redisClient.Close()

	userRepository := postgres.NewUserRepository(postgresPool)
	if envConfig.SeedData {
		seedRepository := postgres.NewSeedRepository(envConfig.Env)
		if err := seedRepository.Seed(ctx, userRepository); err != nil {
			return fmt.Errorf("seed database: %w", err)
		}
	}

	tokenRepository := redis.NewTokenRepository(redisClient)
	userService := appservice.NewUserService(userRepository, auth.NewPasswordService())
	providerIdentityRepository := postgres.NewProviderIdentityRepository(postgresPool)
	oauthService := appservice.NewOAuthService(providerIdentityRepository)
	sessionService := appservice.NewSessionService(
		tokenService, csrfService, tokenRepository, envConfig.JWTConfig,
	)
	authHandler := httptransport.NewAuthHandler(
		tokenRepository, userService, sessionService, tokenService, csrfService, envConfig.JWTConfig,
	)
	userHandler := httptransport.NewUserHandler(userService)
	oauthHandler := oauth2transport.NewGoogleOAuth2(
		envConfig.OAuth2Config, envConfig.JWTConfig, oauthService, sessionService,
	)

	app := httptransport.NewRouter(
		envConfig,
		tokenRepository,
		tokenService,
		userHandler,
		authHandler,
		oauthHandler,
	)

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- httptransport.StartServer(app, envConfig.Port)
	}()

	select {
	case err := <-serverErr:
		if err != nil {
			return fmt.Errorf("serve HTTP: %w", err)
		}
		return errors.New("HTTP server stopped unexpectedly")
	case <-ctx.Done():
		log.Info().Msg("received termination signal, shutting down")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown HTTP server: %w", err)
	}

	if err := <-serverErr; err != nil {
		return fmt.Errorf("serve HTTP: %w", err)
	}
	return nil
}

package config

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/service"
	logger_env "github.com/DarrelA/starter-go-postgresql/internal/infrastructure/logger"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type EnvConfig struct {
	entity.EnvConfig
}

func NewTokenService() service.LoadEnvConfig {
	return &EnvConfig{}
}

func (e *EnvConfig) LoadAppConfig() {
	e.Env = os.Getenv("APP_ENV")
	if e.Env == "" {
		log.Fatal().Msg("APP_ENV not set")
	}

	cwd := logger_env.LogCWD()
	logger_env.ListFiles()
	envBasePath := "configs/"

	// @TODO: Explore test binary compilation with `go test -c`
	// Check if the current working directory contains "\test"
	if strings.Contains(cwd, "\\test") || strings.Contains(cwd, "/test") {
		envBasePath = "../configs/"
	}

	// Construct the full path to the .env file
	envFilePath := filepath.Join(envBasePath, ".env."+e.Env)
	log.Debug().Msgf("loading env file: %s", envFilePath)
	godotenv.Load(envFilePath)

	e.Port = os.Getenv("APP_PORT")
	if e.Port == "" {
		e.Port = "8080" // Default port
	}

	log.Info().Msgf("running in %s env using Port %s", strings.ToUpper(e.Env), e.Port)

	/*
		1.	Embedding `entity.EnvConfig`: In the `config` package, EnvConfig struct embeds the `entity.EnvConfig` struct.
				This means EnvConfig includes all fields and methods of `entity.EnvConfig` by default.
		2.	Pointer Initialization: Since `BaseURLsConfig` in `entity.EnvConfig` is a pointer (`*BaseURLsConfig`),
				you need to initialize it using `&entity.BaseURLsConfig{}` to allocate memory and assign values.
		3.	Correct Typing: The correct way to initialize a struct field from another package is to use the full type name, in this case, `entity.BaseURLsConfig`.
	*/
	e.BaseURLsConfig = &entity.BaseURLsConfig{
		AuthServicePathName: os.Getenv("AUTH_SERVICE_PATHNAME"),
		AuthService: os.Getenv("PROTOCOL") +
			os.Getenv("DOMAIN") + e.Port +
			os.Getenv("AUTH_SERVICE_PATHNAME"),
	}
}

func (e *EnvConfig) LoadLogConfig() {
	logLevel := os.Getenv("LOG_LEVEL")

	// Whichever level is chosen,
	// all logs with a level greater than or equal to that level will be written.
	switch strings.ToLower(logLevel) {
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel) // Level -1
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel) // Level 0
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel) // Level 1
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel) // Level 2
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel) // Level 3
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel) // Level 4
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel) // Level 5
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel) // Level 1
	}

	if _, err := os.Stat("/app/logs"); os.IsNotExist(err) {
		os.Mkdir("/app/logs", 0755)
	}
}

func (e *EnvConfig) LoadDBConfig() {
	e.PostgresDBConfig = &entity.PostgresDBConfig{
		Username:     os.Getenv("POSTGRES_USER"),
		Password:     os.Getenv("POSTGRES_PASSWORD"),
		Host:         os.Getenv("POSTGRES_HOST"),
		Port:         os.Getenv("POSTGRES_PORT"),
		Name:         os.Getenv("POSTGRES_DB"),
		SslMode:      os.Getenv("POSTGRES_SSLMODE"),
		PoolMaxConns: os.Getenv("POSTGRES_POOL_MAX_CONNS"),
	}
}

func (e *EnvConfig) LoadRedisConfig() {
	e.RedisDBConfig = &entity.RedisDBConfig{
		RedisUri: os.Getenv("REDIS_URL"),
	}
}

func (e *EnvConfig) LoadJWTConfig() {
	/*
		Ensure that `JWTConfig` is properly initialized to avoid `nil` pointer dereference errors.
		This check verifies if `JWTConfig` is `nil`, and if it is,
		initializes it with a new instance of `entity.JWTConfig`.
		This is crucial before attempting to access or modify any fields within `JWTConfig`,
		ensuring that subsequent dereference operations are safe.
	*/
	if e.JWTConfig == nil {
		e.JWTConfig = &entity.JWTConfig{}
	}

	e.JWTConfig.Path = os.Getenv("JWT_PATH")
	e.JWTConfig.Domain = os.Getenv("JWT_DOMAIN")
	loadEnvVariableBool("JWT_SECURE", &e.JWTConfig.Secure)
	loadEnvVariableBool("JWT_HTTPONLY", &e.JWTConfig.HttpOnly)
	e.JWTConfig.AccessTokenPrivateKey = os.Getenv("ACCESS_TOKEN_PRIVATE_KEY")
	e.JWTConfig.AccessTokenPublicKey = os.Getenv("ACCESS_TOKEN_PUBLIC_KEY")
	loadEnvVariableDuration("ACCESS_TOKEN_EXPIRED_IN", &e.JWTConfig.AccessTokenExpiredIn)
	loadEnvVariableInt("ACCESS_TOKEN_MAXAGE", &e.JWTConfig.AccessTokenMaxAge)
	e.JWTConfig.RefreshTokenPrivateKey = os.Getenv("REFRESH_TOKEN_PRIVATE_KEY")
	e.JWTConfig.RefreshTokenPublicKey = os.Getenv("REFRESH_TOKEN_PUBLIC_KEY")
	loadEnvVariableDuration("REFRESH_TOKEN_EXPIRED_IN", &e.JWTConfig.RefreshTokenExpiredIn)
	loadEnvVariableInt("REFRESH_TOKEN_MAXAGE", &e.JWTConfig.RefreshTokenMaxAge)
}

func (e *EnvConfig) LoadCORSConfig() {
	e.CORSConfig = &entity.CORSConfig{
		AllowedOrigins: os.Getenv("CORS_ALLOWED_ORIGINS"),
	}
}

func loadEnvVariableInt(envVar string, target *int) {
	valueStr := os.Getenv(envVar)
	if valueStr == "" {
		log.Error().Msgf("%s is not set", envVar)
		return
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Error().Err(err).Msgf("check JWT config: %s", envVar)
		return
	}
	*target = value
}

func loadEnvVariableBool(envVar string, target *bool) {
	valueStr := os.Getenv(envVar)
	if valueStr == "" {
		log.Error().Msgf("%s is not set", envVar)
		return
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		log.Error().Err(err).Msgf("check JWT config: %s", envVar)
		return
	}
	*target = value
}

func loadEnvVariableDuration(envVar string, target *time.Duration) {
	valueStr := os.Getenv(envVar)
	if valueStr == "" {
		log.Error().Msgf("%s is not set", envVar)
		return
	}
	value, err := time.ParseDuration(valueStr)
	if err != nil {
		log.Error().Err(err).Msgf("check JWT config: %s", envVar)
		return
	}
	*target = value
}

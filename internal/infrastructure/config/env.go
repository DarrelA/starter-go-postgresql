package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/service"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	dev  = "dev"
	test = "test"
	prod = "prod"

	errMsgEnvVarNotSet   = "%s is not set or has an invalid value; only 'dev', 'test', or 'prod' are accepted"
	errMsgVarNotSet      = "%s is not set"
	errMsgCheckJWTConfig = "check JWT config: %s"
)

type EnvConfig struct {
	entity.EnvConfig
}

func NewTokenService() service.LoadEnvConfig {
	return &EnvConfig{}
}

func (e *EnvConfig) LoadAppConfig() {
	e.Env = os.Getenv("APP_ENV")
	if e.Env != dev && e.Env != test && e.Env != prod {
		log.Fatal().Msgf(errMsgEnvVarNotSet, "APP_ENV")
	}

	e.Port = checkEmptyEnvVar("APP_PORT")
	log.Info().Msgf("running in %s env using Port %s", strings.ToUpper(e.Env), e.Port)

	authServicePathName := checkEmptyEnvVar("AUTH_SERVICE_PATHNAME")
	protocol := checkEmptyEnvVar("PROTOCOL")
	domain := checkEmptyEnvVar("DOMAIN")

	/*
		1.	Embedding `entity.EnvConfig`: In the `config` package, EnvConfig struct embeds the `entity.EnvConfig` struct.
				This means EnvConfig includes all fields and methods of `entity.EnvConfig` by default.
		2.	Pointer Initialization: Since `BaseURLsConfig` in `entity.EnvConfig` is a pointer (`*BaseURLsConfig`),
				you need to initialize it using `&entity.BaseURLsConfig{}` to allocate memory and assign values.
		3.	Correct Typing: The correct way to initialize a struct field from another package is to use the full type name, in this case, `entity.BaseURLsConfig`.
	*/
	e.BaseURLsConfig = &entity.BaseURLsConfig{
		AuthServicePathName: authServicePathName,
		AuthService:         protocol + domain + e.Port + authServicePathName,
	}
}

func (e *EnvConfig) LoadLogConfig() {
	logLevel := checkEmptyEnvVar("LOG_LEVEL")

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
		Username:     checkEmptyEnvVar("POSTGRES_USER"),
		Password:     checkEmptyEnvVar("POSTGRES_PASSWORD"),
		Host:         checkEmptyEnvVar("POSTGRES_HOST"),
		Port:         checkEmptyEnvVar("POSTGRES_PORT"),
		Name:         checkEmptyEnvVar("POSTGRES_DB"),
		SslMode:      checkEmptyEnvVar("POSTGRES_SSLMODE"),
		PoolMaxConns: checkEmptyEnvVar("POSTGRES_POOL_MAX_CONNS"),
	}
}

func (e *EnvConfig) LoadRedisConfig() {
	e.RedisDBConfig = &entity.RedisDBConfig{
		RedisUri: checkEmptyEnvVar("REDIS_URL"),
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

	e.JWTConfig.Path = checkEmptyEnvVar("JWT_PATH")
	e.JWTConfig.Domain = checkEmptyEnvVar("JWT_DOMAIN")
	loadEnvVariableBool("JWT_SECURE", &e.JWTConfig.Secure)
	loadEnvVariableBool("JWT_HTTPONLY", &e.JWTConfig.HttpOnly)
	e.JWTConfig.AccessTokenPrivateKey = checkEmptyEnvVar("ACCESS_TOKEN_PRIVATE_KEY")
	e.JWTConfig.AccessTokenPublicKey = checkEmptyEnvVar("ACCESS_TOKEN_PUBLIC_KEY")
	loadEnvVariableDuration("ACCESS_TOKEN_EXPIRED_IN", &e.JWTConfig.AccessTokenExpiredIn)
	loadEnvVariableInt("ACCESS_TOKEN_MAXAGE", &e.JWTConfig.AccessTokenMaxAge)
	e.JWTConfig.RefreshTokenPrivateKey = checkEmptyEnvVar("REFRESH_TOKEN_PRIVATE_KEY")
	e.JWTConfig.RefreshTokenPublicKey = checkEmptyEnvVar("REFRESH_TOKEN_PUBLIC_KEY")
	loadEnvVariableDuration("REFRESH_TOKEN_EXPIRED_IN", &e.JWTConfig.RefreshTokenExpiredIn)
	loadEnvVariableInt("REFRESH_TOKEN_MAXAGE", &e.JWTConfig.RefreshTokenMaxAge)
}

func (e *EnvConfig) LoadCORSConfig() {
	e.CORSConfig = &entity.CORSConfig{
		AllowedOrigins: checkEmptyEnvVar("CORS_ALLOWED_ORIGINS"),
	}
}

func checkEmptyEnvVar(envVar string) string {
	valueStr := os.Getenv(envVar)
	if valueStr == "" {
		log.Fatal().Msgf(errMsgVarNotSet, envVar)
	}
	return valueStr
}

func loadEnvVariableInt(envVar string, target *int) {
	valueStr := checkEmptyEnvVar(envVar)
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Error().Err(err).Msgf(errMsgCheckJWTConfig, envVar)
		return
	}
	*target = value
}

func loadEnvVariableBool(envVar string, target *bool) {
	valueStr := checkEmptyEnvVar(envVar)
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		log.Error().Err(err).Msgf(errMsgCheckJWTConfig, envVar)
		return
	}
	*target = value
}

func loadEnvVariableDuration(envVar string, target *time.Duration) {
	valueStr := checkEmptyEnvVar(envVar)
	value, err := time.ParseDuration(valueStr)
	if err != nil {
		log.Error().Err(err).Msgf(errMsgCheckJWTConfig, envVar)
		return
	}
	*target = value
}

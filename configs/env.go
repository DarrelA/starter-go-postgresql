package configs

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	logger_env "github.com/DarrelA/starter-go-postgresql/internal/infrastructure/logger"
	"github.com/DarrelA/starter-go-postgresql/internal/utils/err_rest"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type BaseURLsConfig struct {
	AuthServicePathName string
	AuthService         string
}

type PostgresDBConfig struct {
	Username     string
	Password     string
	Host         string
	Port         string
	Name         string
	SslMode      string
	PoolMaxConns string
}

type RedisDBConfig struct {
	RedisUri string
}

type JWTConfig struct {
	Path                   string
	Domain                 string
	Secure                 bool
	HttpOnly               bool
	AccessTokenPrivateKey  string
	AccessTokenPublicKey   string
	AccessTokenExpiredIn   time.Duration
	AccessTokenMaxAge      int
	RefreshTokenPrivateKey string
	RefreshTokenPublicKey  string
	RefreshTokenExpiredIn  time.Duration
	RefreshTokenMaxAge     int
}

type CORSConfig struct {
	AllowedOrigins string
}

// Define the variables to hold the configuration
var (
	Env          string
	Port         string
	BaseURLs     BaseURLsConfig
	PGDB         PostgresDBConfig
	RedisDB      RedisDBConfig
	JWTSettings  JWTConfig
	CORSSettings CORSConfig
)

func init() {
	loadEnv()
	loadLogSettings()
	loadDBSettings()
	loadRedisSettings()
	loadJWTConfigs()
	loadCORSConfigs()
}

func loadEnv() {
	Env = os.Getenv("APP_ENV")
	if Env == "" {
		log.Fatal().Msg("APP_ENV not set")
	}

	cwd := logger_env.LogCWD()
	logger_env.ListFiles()
	envBasePath := "configs/"

	// @TODO: Explore `integration-test-coverage-issue` branch
	// Check if the current working directory contains "\test"
	if strings.Contains(cwd, "\\test") || strings.Contains(cwd, "/test") {
		envBasePath = "../configs/"
	}

	// Construct the full path to the .env file
	envFilePath := filepath.Join(envBasePath, ".env."+Env)
	log.Debug().Msgf("loading env file: %s", envFilePath)
	godotenv.Load(envFilePath)

	Port = os.Getenv("APP_PORT")
	if Port == "" {
		Port = "8080" // Default port
	}

	log.Info().Msgf("running in %s env using Port %s", strings.ToUpper(Env), Port)

	BaseURLs = BaseURLsConfig{
		AuthServicePathName: os.Getenv("AUTH_SERVICE_PATHNAME"),
		AuthService:         os.Getenv("PROTOCOL") + os.Getenv("DOMAIN") + Port + os.Getenv("AUTH_SERVICE_PATHNAME"),
	}
}

func loadLogSettings() {
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

func loadDBSettings() {
	PGDB = PostgresDBConfig{
		Username:     os.Getenv("POSTGRES_USER"),
		Password:     os.Getenv("POSTGRES_PASSWORD"),
		Host:         os.Getenv("POSTGRES_HOST"),
		Port:         os.Getenv("POSTGRES_PORT"),
		Name:         os.Getenv("POSTGRES_DB"),
		SslMode:      os.Getenv("POSTGRES_SSLMODE"),
		PoolMaxConns: os.Getenv("POSTGRES_POOL_MAX_CONNS"),
	}
}

func loadRedisSettings() {
	RedisDB = RedisDBConfig{
		RedisUri: os.Getenv("REDIS_URL"),
	}
}

func loadJWTConfigs() {
	JWTSettings.Path = os.Getenv("JWT_PATH")
	JWTSettings.Domain = os.Getenv("JWT_DOMAIN")
	loadEnvVariableBool("JWT_SECURE", &JWTSettings.Secure)
	loadEnvVariableBool("JWT_HTTPONLY", &JWTSettings.HttpOnly)
	JWTSettings.AccessTokenPrivateKey = os.Getenv("ACCESS_TOKEN_PRIVATE_KEY")
	JWTSettings.AccessTokenPublicKey = os.Getenv("ACCESS_TOKEN_PUBLIC_KEY")
	loadEnvVariableDuration("ACCESS_TOKEN_EXPIRED_IN", &JWTSettings.AccessTokenExpiredIn)
	loadEnvVariableInt("ACCESS_TOKEN_MAXAGE", &JWTSettings.AccessTokenMaxAge)
	JWTSettings.RefreshTokenPrivateKey = os.Getenv("REFRESH_TOKEN_PRIVATE_KEY")
	JWTSettings.RefreshTokenPublicKey = os.Getenv("REFRESH_TOKEN_PUBLIC_KEY")
	loadEnvVariableDuration("REFRESH_TOKEN_EXPIRED_IN", &JWTSettings.RefreshTokenExpiredIn)
	loadEnvVariableInt("REFRESH_TOKEN_MAXAGE", &JWTSettings.RefreshTokenMaxAge)
}

func loadCORSConfigs() {
	CORSSettings = CORSConfig{
		AllowedOrigins: os.Getenv("CORS_ALLOWED_ORIGINS"),
	}
}

func loadEnvVariableInt(envVar string, target *int) {
	valueStr := os.Getenv(envVar)
	if valueStr == "" {
		errMessage := fmt.Sprintf("%s is not set", envVar)
		err := err_rest.NewInternalServerError(errMessage)
		log.Error().Err(err).Msg("")
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
		errMessage := fmt.Sprintf("%s is not set", envVar)
		err := err_rest.NewInternalServerError(errMessage)
		log.Error().Err(err).Msg("")
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
		errMessage := fmt.Sprintf("%s is not set", envVar)
		err := err_rest.NewInternalServerError(errMessage)
		log.Error().Err(err).Msg("")
		return
	}
	value, err := time.ParseDuration(valueStr)
	if err != nil {
		log.Error().Err(err).Msgf("check JWT config: %s", envVar)
		return
	}
	*target = value
}

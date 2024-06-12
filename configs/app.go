package configs

import (
	"os"
	"strconv"
	"time"

	"github.com/DarrelA/starter-go-postgresql/internal/utils/errors"
	"github.com/joho/godotenv"
)

var Port string

func loadEnv() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	envBasePath := "configs/"

	godotenv.Load(envBasePath + ".env." + env + ".local")
	if env != "test" {
		godotenv.Load(envBasePath + ".env.local")
	}

	godotenv.Load(envBasePath + ".env." + env)
	godotenv.Load()

	Port = os.Getenv("APP_PORT")
	if Port == "" {
		Port = "4040" // Default port
	}
}

type DBConfig struct {
	Username     string
	Password     string
	Host         string
	Port         string
	Name         string
	SslMode      string
	PoolMaxConns string
}

var PGDB DBConfig

func initDBSettings() {
	PGDB = DBConfig{
		Username:     os.Getenv("POSTGRES_USER"),
		Password:     os.Getenv("POSTGRES_PASSWORD"),
		Host:         os.Getenv("POSTGRES_HOST"),
		Port:         os.Getenv("POSTGRES_PORT"),
		Name:         os.Getenv("POSTGRES_NAME"),
		SslMode:      os.Getenv("POSTGRES_SSLMODE"),
		PoolMaxConns: os.Getenv("POSTGRES_POOL_MAX_CONNS"),
	}
}

type RedisConfig struct {
	RedisUri string
}

var RedisDB RedisConfig

func initRedisSettings() {
	RedisDB = RedisConfig{
		RedisUri: os.Getenv("REDIS_URL"),
	}
}

type JWTConfig struct {
	Secret                 string
	Name                   string
	Path                   string
	Domain                 string
	MaxAge                 int
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

// Define a variable to hold the configuration
var JWTSettings JWTConfig

func loadJWTConfigs() {
	MaxAge, err := strconv.Atoi(os.Getenv("JWT_MAXAGE"))
	if err != nil {
		errors.NewInternalServerError("Check JWT Config")
	}

	Secure, err := strconv.ParseBool(os.Getenv("JWT_Secure"))
	if err != nil {
		errors.NewInternalServerError("Check JWT Config")
	}

	HttpOnly, err := strconv.ParseBool(os.Getenv("JWT_HTTPONLY"))
	if err != nil {
		errors.NewInternalServerError("Check JWT Config")
	}

	AccessTokenMaxAge, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_MAXAGE"))
	if err != nil {
		errors.NewInternalServerError("Check JWT Config")
	}

	AccessTokenExpiredIn, err := time.ParseDuration(os.Getenv("ACCESS_TOKEN_EXPIRED_IN"))
	if err != nil {
		panic("Check JWT Config")
	}

	RefreshTokenMaxAge, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_MAXAGE"))
	if err != nil {
		errors.NewInternalServerError("Check JWT Config")
	}

	RefreshTokenExpiredIn, err := time.ParseDuration(os.Getenv("REFRESH_TOKEN_EXPIRED_IN"))
	if err != nil {
		panic("Check JWT Config")
	}

	// Initialize the JWTSettings struct with environment variable values
	JWTSettings = JWTConfig{
		Secret:                 os.Getenv("JWT_SECRET"),
		Name:                   os.Getenv("JWT_NAME"),
		Path:                   os.Getenv("JWT_PATH"),
		Domain:                 os.Getenv("JWT_DOMAIN"),
		MaxAge:                 MaxAge,
		Secure:                 Secure,
		HttpOnly:               HttpOnly,
		AccessTokenPrivateKey:  os.Getenv("ACCESS_TOKEN_PRIVATE_KEY"),
		AccessTokenPublicKey:   os.Getenv("ACCESS_TOKEN_PUBLIC_KEY"),
		AccessTokenExpiredIn:   AccessTokenExpiredIn,
		AccessTokenMaxAge:      AccessTokenMaxAge,
		RefreshTokenPrivateKey: os.Getenv("REFRESH_TOKEN_PRIVATE_KEY"),
		RefreshTokenPublicKey:  os.Getenv("REFRESH_TOKEN_PUBLIC_KEY"),
		RefreshTokenExpiredIn:  RefreshTokenExpiredIn,
		RefreshTokenMaxAge:     RefreshTokenMaxAge,
	}
}

type CORSConfig struct {
	AllowedOrigins string
}

var CORSSettings CORSConfig

func loadCORSConfigs() {
	CORSSettings = CORSConfig{
		AllowedOrigins: os.Getenv("CORS_ALLOWED_ORIGINS"),
	}
}

func init() {
	loadEnv()
	initDBSettings()
	initRedisSettings()
	loadJWTConfigs()
	loadCORSConfigs()
}

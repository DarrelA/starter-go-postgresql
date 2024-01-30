package configs

import (
	"os"
	"strconv"

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
		Username:     os.Getenv("POSTGRES_USERNAME"),
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
	Secret   string
	Name     string
	Path     string
	Domain   string
	MaxAge   int
	Secure   bool
	HttpOnly bool
}

// Define a variable to hold the configuration
var JWTSettings JWTConfig

func loadJWTConfigs() {
	MaxAge, err := strconv.Atoi("JWT_MAXAGE")
	if err != nil {
		errors.NewInternalServerError("Check JWT Config")
	}

	Secure, err := strconv.ParseBool("JWT_SECURE")
	if err != nil {
		errors.NewInternalServerError("Check JWT Config")
	}

	HttpOnly, err := strconv.ParseBool("JWT_HTTPONLY")
	if err != nil {
		errors.NewInternalServerError("Check JWT Config")
	}

	// Initialize the JWTSettings struct with environment variable values
	JWTSettings = JWTConfig{
		Secret:   os.Getenv("JWT_SECRET"),
		Name:     os.Getenv("JWT_NAME"),
		Path:     os.Getenv("JWT_PATH"),
		Domain:   os.Getenv("JWT_DOMAIN"),
		MaxAge:   MaxAge,
		Secure:   Secure,
		HttpOnly: HttpOnly,
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

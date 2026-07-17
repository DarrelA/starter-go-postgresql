package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	"github.com/rs/zerolog"
)

type EnvConfig struct {
	entity.EnvConfig
	LogLevel         string
	SeedData         bool
	MigrateOnStartup bool
}

// LoadPostgresConfig reads only the configuration required by database jobs.
func LoadPostgresConfig() (*entity.PostgresDBConfig, error) {
	reader := envReader{}
	postgresConfig := readPostgresConfig(&reader)
	if err := reader.err(); err != nil {
		return nil, fmt.Errorf("load PostgreSQL configuration: %w", err)
	}
	return postgresConfig, nil
}

// LoadEnvConfig reads and validates all runtime configuration.
func LoadEnvConfig() (*EnvConfig, error) {
	reader := envReader{}

	env := reader.required("APP_ENV")
	if env != "dev" && env != "test" && env != "prod" {
		reader.add("APP_ENV must be one of dev, test, or prod")
	}

	port := reader.positiveInt("APP_PORT")
	if port > 65535 {
		reader.add("APP_PORT must be at most 65535")
	}

	protocol := reader.required("PROTOCOL")
	if protocol != "http://" && protocol != "https://" {
		reader.add("PROTOCOL must be http:// or https://")
	}
	domain := reader.required("DOMAIN")
	authPath := reader.required("AUTH_SERVICE_PATHNAME")
	if authPath != "" && !strings.HasPrefix(authPath, "/") {
		reader.add("AUTH_SERVICE_PATHNAME must start with /")
	}

	logLevel := strings.ToLower(reader.required("LOG_LEVEL"))
	if _, err := zerolog.ParseLevel(logLevel); err != nil {
		reader.add("LOG_LEVEL is invalid")
	}

	jwtPath := reader.required("JWT_PATH")
	if jwtPath != "" && !strings.HasPrefix(jwtPath, "/") {
		reader.add("JWT_PATH must start with /")
	}
	csrfSecret := reader.required("CSRF_SECRET")
	if len(csrfSecret) < 32 {
		reader.add("CSRF_SECRET must contain at least 32 characters")
	}
	jwtHTTPOnly := reader.boolean("JWT_HTTPONLY")
	if !jwtHTTPOnly {
		reader.add("JWT_HTTPONLY must be true")
	}

	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	if (googleClientID == "") != (googleClientSecret == "") {
		reader.add("GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET must be set together")
	}

	config := &EnvConfig{
		EnvConfig: entity.EnvConfig{
			Env:  env,
			Port: strconv.Itoa(port),
			BaseURLsConfig: &entity.BaseURLsConfig{
				AuthServicePathName: authPath,
				AuthService:         protocol + domain + ":" + strconv.Itoa(port) + authPath,
			},
			PostgresDBConfig: readPostgresConfig(&reader),
			RedisDBConfig:    &entity.RedisDBConfig{RedisUri: reader.required("REDIS_URL")},
			JWTConfig: &entity.JWTConfig{
				Path:                   jwtPath,
				Domain:                 reader.required("JWT_DOMAIN"),
				Secure:                 reader.boolean("JWT_SECURE"),
				HttpOnly:               jwtHTTPOnly,
				AccessTokenPrivateKey:  reader.required("ACCESS_TOKEN_PRIVATE_KEY"),
				AccessTokenPublicKey:   reader.required("ACCESS_TOKEN_PUBLIC_KEY"),
				AccessTokenExpiredIn:   reader.duration("ACCESS_TOKEN_EXPIRED_IN"),
				AccessTokenMaxAge:      reader.positiveInt("ACCESS_TOKEN_MAXAGE"),
				RefreshTokenPrivateKey: reader.required("REFRESH_TOKEN_PRIVATE_KEY"),
				RefreshTokenPublicKey:  reader.required("REFRESH_TOKEN_PUBLIC_KEY"),
				RefreshTokenExpiredIn:  reader.duration("REFRESH_TOKEN_EXPIRED_IN"),
				RefreshTokenMaxAge:     reader.positiveInt("REFRESH_TOKEN_MAXAGE"),
			},
			CSRFConfig: &entity.CSRFConfig{Secret: csrfSecret},
			CORSConfig: &entity.CORSConfig{AllowedOrigins: reader.required("CORS_ALLOWED_ORIGINS")},
			OAuth2Config: &entity.OAuth2Config{
				GoogleRedirectURL:  protocol + domain + ":" + strconv.Itoa(port) + "/auth/google_callback",
				GoogleClientID:     googleClientID,
				GoogleClientSecret: googleClientSecret,
				Scopes: []string{
					"https://www.googleapis.com/auth/userinfo.email",
					"https://www.googleapis.com/auth/userinfo.profile",
				},
			},
		},
		LogLevel:         logLevel,
		SeedData:         reader.boolean("SEED_DATA"),
		MigrateOnStartup: reader.optionalBoolean("MIGRATE_ON_STARTUP", env != "prod"),
	}

	if err := reader.err(); err != nil {
		return nil, fmt.Errorf("load configuration: %w", err)
	}
	return config, nil
}

func readPostgresConfig(reader *envReader) *entity.PostgresDBConfig {
	return &entity.PostgresDBConfig{
		Username:     reader.required("POSTGRES_USER"),
		Password:     reader.required("POSTGRES_PASSWORD"),
		Host:         reader.required("POSTGRES_HOST"),
		Port:         reader.positiveIntString("POSTGRES_PORT"),
		Name:         reader.required("POSTGRES_DB"),
		SslMode:      reader.required("POSTGRES_SSLMODE"),
		PoolMaxConns: reader.positiveIntString("POSTGRES_POOL_MAX_CONNS"),
	}
}

// ConfigureLogging applies a previously validated log level.
func ConfigureLogging(level string) error {
	parsedLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		return fmt.Errorf("configure logging: %w", err)
	}
	zerolog.SetGlobalLevel(parsedLevel)
	return nil
}

type envReader struct {
	errors []error
}

func (r *envReader) add(message string) {
	r.errors = append(r.errors, errors.New(message))
}

func (r *envReader) required(name string) string {
	value := os.Getenv(name)
	if value == "" {
		r.add(name + " is required")
	}
	return value
}

func (r *envReader) positiveInt(name string) int {
	value := r.required(name)
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		r.add(name + " must be a positive integer")
		return 0
	}
	return parsed
}

func (r *envReader) positiveIntString(name string) string {
	value := r.required(name)
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		r.add(name + " must be a positive integer")
	}
	return value
}

func (r *envReader) boolean(name string) bool {
	value := r.required(name)
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		r.add(name + " must be true or false")
	}
	return parsed
}

func (r *envReader) optionalBoolean(name string, fallback bool) bool {
	value := os.Getenv(name)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		r.add(name + " must be true or false")
	}
	return parsed
}

func (r *envReader) duration(name string) time.Duration {
	value := r.required(name)
	parsed, err := time.ParseDuration(value)
	if err != nil || parsed <= 0 {
		r.add(name + " must be a positive duration")
		return 0
	}
	return parsed
}

func (r *envReader) err() error {
	return errors.Join(r.errors...)
}

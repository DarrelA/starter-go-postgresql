package config

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestLoadAppConfig(t *testing.T) {
	tests := []struct {
		name        string
		env         string
		expectedEnv string
	}{
		{name: "ValidDevEnv", env: "dev", expectedEnv: "dev"},
		{name: "EmptyEnv", env: "", expectedEnv: "test"},
		{name: "InvalidEnv", env: "invalid_env", expectedEnv: "test"},
	}

	e := &EnvConfig{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("APP_ENV", tt.env)
			defer os.Unsetenv("APP_ENV")
			e.LoadAppConfig()
			if e.Env != tt.expectedEnv {
				t.Errorf("expected e.Env to be '%s', got '%s'", tt.expectedEnv, e.Env)
			}
		})
	}

	os.Setenv("APP_PORT", "8081")
	os.Setenv("AUTH_SERVICE_PATHNAME", "/auth")
	os.Setenv("PROTOCOL", "http://")
	os.Setenv("DOMAIN", "localhost")

	defer os.Unsetenv("APP_PORT")
	defer os.Unsetenv("AUTH_SERVICE_PATHNAME")
	defer os.Unsetenv("PROTOCOL")
	defer os.Unsetenv("DOMAIN")

	e.LoadAppConfig()

	if e.Port != "8081" {
		t.Errorf("expected Port to be '8081', got '%s'", e.Port)
	}

	if e.BaseURLsConfig == nil {
		t.Errorf("BaseURLsConfig is nil")
	}

	if e.BaseURLsConfig.AuthService != "http://localhost:8081/auth" {
		t.Errorf("expected AuthService to be 'http://localhost:8081/auth', got '%s'", e.BaseURLsConfig.AuthService)
	}
}

func TestLoadLogConfig(t *testing.T) {
	tests := []struct {
		name             string
		envConfig        EnvConfig
		logLevel         string
		expectedLogLevel zerolog.Level
	}{
		{name: "TraceLogLevel",
			envConfig: EnvConfig{EnvConfig: entity.EnvConfig{Env: "prod"}},
			logLevel:  "trace", expectedLogLevel: zerolog.TraceLevel},

		{name: "InfoLogLevel",
			envConfig: EnvConfig{EnvConfig: entity.EnvConfig{Env: "prod"}},
			logLevel:  "info", expectedLogLevel: zerolog.InfoLevel},

		{name: "WarnLogLevel",
			envConfig: EnvConfig{EnvConfig: entity.EnvConfig{Env: "prod"}},
			logLevel:  "warn", expectedLogLevel: zerolog.WarnLevel},

		{name: "ErrorLogLevel",
			envConfig: EnvConfig{EnvConfig: entity.EnvConfig{Env: "prod"}},
			logLevel:  "error", expectedLogLevel: zerolog.ErrorLevel},

		{name: "FatalLogLevel",
			envConfig: EnvConfig{EnvConfig: entity.EnvConfig{Env: "prod"}},
			logLevel:  "fatal", expectedLogLevel: zerolog.FatalLevel},

		{name: "ValidCaseSensitiveLogLevel",
			envConfig: EnvConfig{EnvConfig: entity.EnvConfig{Env: "prod"}},
			logLevel:  "paNiC", expectedLogLevel: zerolog.PanicLevel},

		{name: "EmptyLogLevel",
			envConfig: EnvConfig{EnvConfig: entity.EnvConfig{Env: "prod"}},
			logLevel:  "", expectedLogLevel: zerolog.InfoLevel},

		{name: "InvalidLogLevel",
			envConfig: EnvConfig{EnvConfig: entity.EnvConfig{Env: "dev"}},
			logLevel:  "invalid_logLeVeL", expectedLogLevel: zerolog.DebugLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("LOG_LEVEL", tt.logLevel)
			defer os.Unsetenv("LOG_LEVEL")
			e := tt.envConfig
			e.LoadLogConfig()
			if zerolog.GlobalLevel() != tt.expectedLogLevel {
				t.Errorf("expected zerolog.GlobalLevel() to be '%s', got '%s'", tt.expectedLogLevel, zerolog.GlobalLevel())
			}
		})
	}
}

func TestLoadDBConfig(t *testing.T) {
	e := &EnvConfig{}

	os.Setenv("POSTGRES_USER", "Only checkEmptyEnvVar validation")
	os.Setenv("POSTGRES_PASSWORD", "Only checkEmptyEnvVar validation")
	os.Setenv("POSTGRES_HOST", "Only checkEmptyEnvVar validation")
	os.Setenv("POSTGRES_PORT", "Only checkEmptyEnvVar validation")
	os.Setenv("POSTGRES_DB", "Only checkEmptyEnvVar validation")
	os.Setenv("POSTGRES_SSLMODE", "Only checkEmptyEnvVar validation")
	os.Setenv("POSTGRES_POOL_MAX_CONNS", "Only checkEmptyEnvVar validation")

	defer os.Unsetenv("POSTGRES_USER")
	defer os.Unsetenv("POSTGRES_PASSWORD")
	defer os.Unsetenv("POSTGRES_HOST")
	defer os.Unsetenv("POSTGRES_PORT")
	defer os.Unsetenv("POSTGRES_DB")
	defer os.Unsetenv("POSTGRES_SSLMODE")
	defer os.Unsetenv("POSTGRES_POOL_MAX_CONNS")

	e.LoadDBConfig()

	if e.PostgresDBConfig == nil {
		t.Errorf("PostgresDBConfig is nil")
	}

	if e.PostgresDBConfig.Username != "Only checkEmptyEnvVar validation" {
		t.Errorf("expected Username to be 'Only checkEmptyEnvVar validation', got '%s'", e.PostgresDBConfig.Username)
	}
	if e.PostgresDBConfig.Password != "Only checkEmptyEnvVar validation" {
		t.Errorf("expected Password to be 'Only checkEmptyEnvVar validation', got '%s'", e.PostgresDBConfig.Password)
	}
	if e.PostgresDBConfig.Host != "Only checkEmptyEnvVar validation" {
		t.Errorf("expected Host to be 'Only checkEmptyEnvVar validation', got '%s'", e.PostgresDBConfig.Host)
	}
	if e.PostgresDBConfig.Port != "Only checkEmptyEnvVar validation" {
		t.Errorf("expected Port to be 'Only checkEmptyEnvVar validation', got '%s'", e.PostgresDBConfig.Port)
	}
	if e.PostgresDBConfig.Name != "Only checkEmptyEnvVar validation" {
		t.Errorf("expected Name to be 'Only checkEmptyEnvVar validation', got '%s'", e.PostgresDBConfig.Name)
	}
	if e.PostgresDBConfig.SslMode != "Only checkEmptyEnvVar validation" {
		t.Errorf("expected SslMode to be 'Only checkEmptyEnvVar validation', got '%s'", e.PostgresDBConfig.SslMode)
	}
	if e.PostgresDBConfig.PoolMaxConns != "Only checkEmptyEnvVar validation" {
		t.Errorf("expected PoolMaxConns to be 'Only checkEmptyEnvVar validation', got '%s'", e.PostgresDBConfig.PoolMaxConns)
	}
}

func TestLoadRedisConfig(t *testing.T) {
	e := &EnvConfig{}
	os.Setenv("REDIS_URL", "Only checkEmptyEnvVar validation")
	defer os.Unsetenv("REDIS_URL")
	e.LoadRedisConfig()

	if e.RedisDBConfig == nil {
		t.Errorf("RedisDBConfig is nil")
	}

	if e.RedisDBConfig.RedisUri != "Only checkEmptyEnvVar validation" {
		t.Errorf("expected RedisUri to be 'Only checkEmptyEnvVar validation', got '%s'", e.RedisDBConfig.RedisUri)
	}
}

func TestLoadJWTConfig(t *testing.T) {
	expiredIn := 1200 * time.Minute
	maxAge := 600

	e := &EnvConfig{}
	os.Setenv("JWT_PATH", "Only checkEmptyEnvVar validation")
	os.Setenv("JWT_DOMAIN", "Only checkEmptyEnvVar validation")
	os.Setenv("JWT_SECURE", "trUe")
	os.Setenv("JWT_HTTPONLY", "trUe")
	os.Setenv("ACCESS_TOKEN_PRIVATE_KEY", "Only checkEmptyEnvVar validation")
	os.Setenv("ACCESS_TOKEN_PUBLIC_KEY", "Only checkEmptyEnvVar validation")
	os.Setenv("ACCESS_TOKEN_EXPIRED_IN", "1200m")
	os.Setenv("ACCESS_TOKEN_MAXAGE", "600")
	os.Setenv("REFRESH_TOKEN_PRIVATE_KEY", "Only checkEmptyEnvVar validation")
	os.Setenv("REFRESH_TOKEN_PUBLIC_KEY", "Only checkEmptyEnvVar validation")
	os.Setenv("REFRESH_TOKEN_EXPIRED_IN", "1200m")
	os.Setenv("REFRESH_TOKEN_MAXAGE", "600")

	defer os.Unsetenv("JWT_PATH")
	defer os.Unsetenv("JWT_DOMAIN")
	defer os.Unsetenv("JWT_SECURE")
	defer os.Unsetenv("JWT_HTTPONLY")
	defer os.Unsetenv("ACCESS_TOKEN_PRIVATE_KEY")
	defer os.Unsetenv("ACCESS_TOKEN_PUBLIC_KEY")
	defer os.Unsetenv("ACCESS_TOKEN_EXPIRED_IN")
	defer os.Unsetenv("ACCESS_TOKEN_MAXAGE")
	defer os.Unsetenv("REFRESH_TOKEN_PRIVATE_KEY")
	defer os.Unsetenv("REFRESH_TOKEN_PUBLIC_KEY")
	defer os.Unsetenv("REFRESH_TOKEN_EXPIRED_IN")
	defer os.Unsetenv("REFRESH_TOKEN_MAXAGE")

	e.LoadJWTConfig()

	if e.JWTConfig == nil {
		t.Errorf("JWTConfig is nil")
	}

	if e.JWTConfig.Path != "Only checkEmptyEnvVar validation" {
		t.Errorf("expected Path to be 'Only checkEmptyEnvVar validation', got '%s'", e.JWTConfig.Path)
	}
	if e.JWTConfig.Domain != "Only checkEmptyEnvVar validation" {
		t.Errorf("expected Domain to be 'Only checkEmptyEnvVar validation', got '%s'", e.JWTConfig.Domain)
	}
	if e.JWTConfig.Secure != true {
		t.Errorf("expected Secure to be 'true', got '%t'", e.JWTConfig.Secure)
	}
	if e.JWTConfig.HttpOnly != true {
		t.Errorf("expected HttpOnly to be 'true', got '%t'", e.JWTConfig.HttpOnly)
	}
	if e.JWTConfig.AccessTokenPrivateKey != "Only checkEmptyEnvVar validation" {
		t.Errorf("expected AccessTokenPrivateKey to be 'Only checkEmptyEnvVar validation', got '%s'",
			e.JWTConfig.AccessTokenPrivateKey)
	}
	if e.JWTConfig.AccessTokenPublicKey != "Only checkEmptyEnvVar validation" {
		t.Errorf("expected AccessTokenPublicKey to be 'Only checkEmptyEnvVar validation', got '%s'",
			e.JWTConfig.AccessTokenPublicKey)
	}
	if e.JWTConfig.AccessTokenExpiredIn != expiredIn {
		t.Errorf("expected AccessTokenExpiredIn to be '%s', got '%s'",
			expiredIn.String(), e.JWTConfig.AccessTokenExpiredIn.String())
	}
	if e.JWTConfig.AccessTokenMaxAge != maxAge {
		t.Errorf("expected AccessTokenMaxAge to be '%d', got '%d'", maxAge, e.JWTConfig.AccessTokenMaxAge)
	}
	if e.JWTConfig.RefreshTokenPrivateKey != "Only checkEmptyEnvVar validation" {
		t.Errorf("expected RefreshTokenPrivateKey to be 'Only checkEmptyEnvVar validation', got '%s'", e.JWTConfig.RefreshTokenPrivateKey)
	}
	if e.JWTConfig.RefreshTokenPublicKey != "Only checkEmptyEnvVar validation" {
		t.Errorf("expected RefreshTokenPublicKey to be 'Only checkEmptyEnvVar validation', got '%s'", e.JWTConfig.RefreshTokenPublicKey)
	}
	if e.JWTConfig.RefreshTokenExpiredIn != expiredIn {
		t.Errorf("expected AccessTokenExpiredIn to be '%s', got '%s'",
			expiredIn.String(), e.JWTConfig.RefreshTokenExpiredIn.String())
	}
	if e.JWTConfig.RefreshTokenMaxAge != maxAge {
		t.Errorf("expected RefreshTokenMaxAge to be '%d', got '%d'", maxAge, e.JWTConfig.RefreshTokenMaxAge)
	}
}

func TestLoadCORSConfig(t *testing.T) {
	e := &EnvConfig{}
	os.Setenv("CORS_ALLOWED_ORIGINS", "Only checkEmptyEnvVar validation")
	defer os.Unsetenv("CORS_ALLOWED_ORIGINS")
	e.LoadCORSConfig()

	if e.CORSConfig == nil {
		t.Errorf("CORSConfig is nil")
	}

	if e.CORSConfig.AllowedOrigins != "Only checkEmptyEnvVar validation" {
		t.Errorf("expected AllowedOrigins to be 'Only checkEmptyEnvVar validation', got '%s'", e.CORSConfig.AllowedOrigins)
	}
}

func TestCheckEmptyEnvVar(t *testing.T) {
	// Create a buffer to capture stdout
	var buf bytes.Buffer

	// Set zerolog to write to the buffer
	log.Logger = log.Output(&buf)
	zerolog.SetGlobalLevel(zerolog.ErrorLevel)

	emptyField := ""
	checkEmptyEnvVar(emptyField)
	logOutput := buf.String()

	// Check if the string contains the expected error message
	expectedMessage := fmt.Sprintf(errMsgVarNotSet, emptyField)
	if !strings.Contains(logOutput, expectedMessage) {
		t.Errorf("Expected error message not found in log output: %s", logOutput)
	}

	// Clean up by resetting the logger
	log.Logger = log.Output(os.Stdout)
}

func TestLoadEnvVariableInt(t *testing.T) {
	var buf bytes.Buffer
	log.Logger = log.Output(&buf)
	zerolog.SetGlobalLevel(zerolog.ErrorLevel)

	negativeNumberInStr := "-1"
	var target int
	loadEnvVariableInt(negativeNumberInStr, &target)
	logOutput := buf.String()

	expectedMessage := fmt.Sprintf(errMsgVarNotSet, negativeNumberInStr)
	if !strings.Contains(logOutput, expectedMessage) {
		t.Errorf("Expected error message not found in log output: %s", logOutput)
	}

	log.Logger = log.Output(os.Stdout)
}

func TestLoadEnvVariableBool(t *testing.T) {
	var buf bytes.Buffer
	log.Logger = log.Output(&buf)
	zerolog.SetGlobalLevel(zerolog.ErrorLevel)

	notBool := "fake"
	var target bool
	loadEnvVariableBool(notBool, &target)
	logOutput := buf.String()

	expectedMessage := fmt.Sprintf(errMsgVarNotSet, notBool)
	if !strings.Contains(logOutput, expectedMessage) {
		t.Errorf("Expected error message not found in log output: %s", logOutput)
	}

	log.Logger = log.Output(os.Stdout)
}

func TestLoadEnvVariableDuration(t *testing.T) {
	var buf bytes.Buffer
	log.Logger = log.Output(&buf)
	zerolog.SetGlobalLevel(zerolog.ErrorLevel)

	invalidUnit := "1IU"
	var target time.Duration
	loadEnvVariableDuration(invalidUnit, &target)
	logOutput := buf.String()

	expectedMessage := fmt.Sprintf(errMsgVarNotSet, invalidUnit)
	if !strings.Contains(logOutput, expectedMessage) {
		t.Errorf("Expected error message not found in log output: %s", logOutput)
	}

	log.Logger = log.Output(os.Stdout)
}

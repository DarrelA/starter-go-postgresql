package config

import (
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
)

func TestLoadEnvConfig(t *testing.T) {
	setValidEnvironment(t)

	config, err := LoadEnvConfig()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if config.Env != "test" || config.Port != "8080" || !config.SeedData || !config.MigrateOnStartup {
		t.Fatalf("unexpected application config: %#v", config)
	}
	if config.BaseURLsConfig.AuthService != "http://localhost:8080/auth/api/v1/users" {
		t.Fatalf("unexpected auth service URL: %s", config.BaseURLsConfig.AuthService)
	}
	if config.JWTConfig.AccessTokenExpiredIn != time.Minute {
		t.Fatalf("unexpected access-token duration: %s", config.JWTConfig.AccessTokenExpiredIn)
	}
}

func TestLoadPostgresConfigDoesNotRequireApplicationConfiguration(t *testing.T) {
	setValidEnvironment(t)
	for _, key := range []string{"APP_PORT", "REDIS_URL", "ACCESS_TOKEN_PRIVATE_KEY", "GOOGLE_CLIENT_ID"} {
		t.Setenv(key, "")
	}
	config, err := LoadPostgresConfig()
	if err != nil {
		t.Fatalf("load migration configuration: %v", err)
	}
	if config.Host != "postgres" || config.Name != "starter" {
		t.Fatalf("unexpected PostgreSQL configuration: %#v", config)
	}
}

func TestLoadEnvConfigAllowsStartupMigrationOverride(t *testing.T) {
	setValidEnvironment(t)
	t.Setenv("MIGRATE_ON_STARTUP", "false")
	config, err := LoadEnvConfig()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if config.MigrateOnStartup {
		t.Fatal("expected startup migration to be disabled")
	}
}

func TestLoadEnvConfigDisablesProductionStartupMigrationByDefault(t *testing.T) {
	setValidEnvironment(t)
	t.Setenv("APP_ENV", "prod")
	t.Setenv("MIGRATE_ON_STARTUP", "")
	config, err := LoadEnvConfig()
	if err != nil {
		t.Fatalf("load production config: %v", err)
	}
	if config.MigrateOnStartup {
		t.Fatal("production startup migration should default to disabled")
	}
}

func TestLoadEnvConfigRejectsInvalidValues(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		value   string
		message string
	}{
		{name: "missing value", key: "POSTGRES_HOST", value: "", message: "POSTGRES_HOST is required"},
		{name: "invalid environment", key: "APP_ENV", value: "staging", message: "APP_ENV must be one of"},
		{name: "invalid boolean", key: "JWT_SECURE", value: "sometimes", message: "JWT_SECURE must be true or false"},
		{name: "invalid duration", key: "ACCESS_TOKEN_EXPIRED_IN", value: "0s", message: "ACCESS_TOKEN_EXPIRED_IN must be a positive duration"},
		{name: "invalid port", key: "APP_PORT", value: "70000", message: "APP_PORT must be at most 65535"},
		{name: "invalid migration switch", key: "MIGRATE_ON_STARTUP", value: "sometimes", message: "MIGRATE_ON_STARTUP must be true or false"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			setValidEnvironment(t)
			t.Setenv(test.key, test.value)
			_, err := LoadEnvConfig()
			if err == nil || !strings.Contains(err.Error(), test.message) {
				t.Fatalf("expected %q, got %v", test.message, err)
			}
		})
	}
}

func TestLoadEnvConfigRequiresOAuthPair(t *testing.T) {
	setValidEnvironment(t)
	t.Setenv("GOOGLE_CLIENT_ID", "client")
	t.Setenv("GOOGLE_CLIENT_SECRET", "")

	_, err := LoadEnvConfig()
	if err == nil || !strings.Contains(err.Error(), "must be set together") {
		t.Fatalf("expected OAuth pair validation, got %v", err)
	}
}

func TestConfigureLogging(t *testing.T) {
	previous := zerolog.GlobalLevel()
	t.Cleanup(func() { zerolog.SetGlobalLevel(previous) })

	if err := ConfigureLogging("warn"); err != nil {
		t.Fatalf("configure logging: %v", err)
	}
	if zerolog.GlobalLevel() != zerolog.WarnLevel {
		t.Fatalf("expected warn level, got %s", zerolog.GlobalLevel())
	}
	if err := ConfigureLogging("verbose"); err == nil {
		t.Fatal("expected invalid log-level error")
	}
}

func setValidEnvironment(t *testing.T) {
	t.Helper()
	values := map[string]string{
		"APP_ENV": "test", "APP_PORT": "8080", "AUTH_SERVICE_PATHNAME": "/auth/api/v1/users",
		"PROTOCOL": "http://", "DOMAIN": "localhost", "LOG_LEVEL": "debug", "SEED_DATA": "true",
		"POSTGRES_USER": "app", "POSTGRES_PASSWORD": "secret", "POSTGRES_HOST": "postgres",
		"POSTGRES_PORT": "5432", "POSTGRES_DB": "starter", "POSTGRES_SSLMODE": "disable",
		"POSTGRES_POOL_MAX_CONNS": "10", "REDIS_URL": "redis:6379", "JWT_PATH": "/",
		"JWT_DOMAIN": "localhost", "JWT_SECURE": "false", "JWT_HTTPONLY": "true",
		"CSRF_SECRET":              "01234567890123456789012345678901",
		"ACCESS_TOKEN_PRIVATE_KEY": "private", "ACCESS_TOKEN_PUBLIC_KEY": "public",
		"ACCESS_TOKEN_EXPIRED_IN": "1m", "ACCESS_TOKEN_MAXAGE": "1",
		"REFRESH_TOKEN_PRIVATE_KEY": "private", "REFRESH_TOKEN_PUBLIC_KEY": "public",
		"REFRESH_TOKEN_EXPIRED_IN": "3m", "REFRESH_TOKEN_MAXAGE": "3",
		"CORS_ALLOWED_ORIGINS": "http://localhost:3030", "GOOGLE_CLIENT_ID": "", "GOOGLE_CLIENT_SECRET": "",
		"MIGRATE_ON_STARTUP": "",
	}
	for key, value := range values {
		t.Setenv(key, value)
	}
}

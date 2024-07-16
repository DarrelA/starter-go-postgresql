package entity

import "time"

type (
	EnvConfig struct {
		Env              string
		Port             string
		BaseURLsConfig   *BaseURLsConfig
		PostgresDBConfig *PostgresDBConfig
		RedisDBConfig    *RedisDBConfig
		JWTConfig        *JWTConfig
		CORSConfig       *CORSConfig
	}

	BaseURLsConfig struct {
		AuthServicePathName string
		AuthService         string
	}

	PostgresDBConfig struct {
		Username     string
		Password     string
		Host         string
		Port         string
		Name         string
		SslMode      string
		PoolMaxConns string
	}

	RedisDBConfig struct {
		RedisUri string
	}

	JWTConfig struct {
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

	CORSConfig struct {
		AllowedOrigins string
	}
)

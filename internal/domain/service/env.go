package service

type LoadEnvConfig interface {
	LoadAppConfig()
	LoadLogConfig()
	LoadDBConfig()
	LoadRedisConfig()
	LoadJWTConfig()
	LoadCORSConfig()
}

package config

/*
Layer Responsibility: Configuration management is often handled at the application layer or infrastructure layer since it involves setting up the environment in which the application runs rather than defining business rules or domain logic.
*/
type LoadEnvConfig interface {
	LoadAppConfig()
	LoadLogConfig()
	LoadDBConfig()
	LoadRedisConfig()
	LoadJWTConfig()
	LoadCORSConfig()
	LoadOAuth2Config()
}

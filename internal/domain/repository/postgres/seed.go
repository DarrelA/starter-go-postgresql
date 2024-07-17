package repository

type PostgresSeedRepository interface {
	Seed(ur PostgresUserRepository)
}

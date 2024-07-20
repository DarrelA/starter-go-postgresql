package postgres

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"strings"

	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	restDomainErr "github.com/DarrelA/starter-go-postgresql/internal/domain/error/transport/http"
	repo "github.com/DarrelA/starter-go-postgresql/internal/domain/repository/postgres"
	password "github.com/DarrelA/starter-go-postgresql/internal/infrastructure/bcrypt"
	logger_env "github.com/DarrelA/starter-go-postgresql/internal/infrastructure/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type PostgresSeedRepository struct {
	dbpool      *pgxpool.Pool
	env         string
	envBasePath string
}

func NewSeedRepository(dbpool *pgxpool.Pool, env string) repo.PostgresSeedRepository {
	envBasePath := "/root/deployment/build"

	cwd := logger_env.LogCWD()
	logger_env.ListFiles()

	// @TODO: Explore test binary compilation with `go test -c`
	// Check if the current working directory contains "\test"
	if strings.Contains(cwd, "\\test") || strings.Contains(cwd, "/test") {
		envBasePath = "../build/"
	}

	ctx := context.Background()

	// Create `users` table in Postgres
	sqlData, err := os.ReadFile(envBasePath + "/sql" + "/schema.user.sql")
	if err != nil {
		log.Error().Err(err).Msgf("unable to read %s/sql/schema.user.sql", envBasePath)
	}
	_, err = dbpool.Exec(ctx, string(sqlData))

	if err != nil {
		log.Error().Err(err).Msg("unable to execute schema.user.sql")
	}

	log.Info().Msg("successfully created extension and table")

	return &PostgresSeedRepository{dbpool, env, envBasePath}
}

func (sr PostgresSeedRepository) Seed(ur repo.PostgresUserRepository) {
	currentEnv := sr.env
	switch currentEnv {
	case "dev":
		saveMultipleUsers(currentEnv, sr.envBasePath, ur)
	case "test":
		saveMultipleUsers(currentEnv, sr.envBasePath, ur)
	default:
		log.Info().Msgf("[%s] env will NOT be seeded with data", currentEnv)
	}
}

func saveMultipleUsers(
	currentEnv string,
	envBasePath string,
	ur repo.PostgresUserRepository,
) *restDomainErr.RestErr {
	userJsonFilePath := "/seed.user." + currentEnv + ".json"
	uu, err := loadUsersFromJsonFile(envBasePath + "/json" + userJsonFilePath)
	if err != nil {
		log.Error().Err(err).Msgf("unable to load [%s]", userJsonFilePath)
	}

	// Verify data in users table by checking for returned errors
	hasData := ur.GetUserByEmail(uu[0])
	if hasData == nil {
		log.Info().Msgf("[%s] env already has seeded data", currentEnv)
		return nil
	}

	for i, u := range uu {
		hashedPassword, _ := password.HashPassword(u.Password)
		uu[i].Password = hashedPassword
		ur.SaveUser(uu[i])
	}

	log.Info().Msgf("successfully seeded data in [%s] env", currentEnv)
	return nil
}

func loadUsersFromJsonFile(filePath string) ([]*entity.User, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var users []*entity.User
	if err := json.Unmarshal(byteValue, &users); err != nil {
		return nil, err
	}

	return users, nil
}

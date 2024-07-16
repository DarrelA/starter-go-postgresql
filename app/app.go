package app

import (
	"context"
	"os"
	"strings"

	"github.com/DarrelA/starter-go-postgresql/db"
	redisDb "github.com/DarrelA/starter-go-postgresql/db/redis"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	logger_env "github.com/DarrelA/starter-go-postgresql/internal/infrastructure/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

var (
	inMemoryDbInstance db.InMemoryDB
)

func CreateRedisConnection(config *entity.RedisDBConfig) db.InMemoryDB {
	if inMemoryDbInstance == nil {
		inMemoryDbInstance = redisDb.NewDB(config)
		inMemoryDbInstance.Connect(config)
	}
	return inMemoryDbInstance
}

func SeedDatabase(dbpool *pgxpool.Pool) {
	envBasePath := "/root/build"
	// currentEnv := configs.Env

	cwd := logger_env.LogCWD()
	logger_env.ListFiles()

	// @TODO: Explore `integration-test-coverage-issue` branch
	// Check if the current working directory contains "\test"
	if strings.Contains(cwd, "\\test") || strings.Contains(cwd, "/test") {
		envBasePath = "../build/"
	}

	ctx := context.Background()

	err := executeSQLFile(ctx, dbpool, envBasePath+"/sql"+"/schema.user.sql")
	if err != nil {
		log.Fatal().Err(err).Msg("unable to execute schema.user.sql")
	}

	log.Info().Msg("successfully created extension and table")

	/*
		switch currentEnv {
		case "dev":
			saveMultipleUsers(currentEnv, envBasePath)
		case "test":
			saveMultipleUsers(currentEnv, envBasePath)
		default:
			log.Info().Msgf("[%s] env will NOT be seeded with data", currentEnv)
		}
	*/
}

func executeSQLFile(ctx context.Context, db *pgxpool.Pool, filePath string) error {
	sqlData, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	// Execute the SQL file
	_, err = db.Exec(ctx, string(sqlData))
	return err
}

/*
		func saveMultipleUsers(currentEnv string, envBasePath string) *err_rest.RestErr {
	userJsonFilePath := "/seed.user." + currentEnv + ".json"
	uu, err := loadUsersFromJsonFile(envBasePath + "/json" + userJsonFilePath)
	if err != nil {
		log.Error().Err(err).Msgf("unable to load [%s]", userJsonFilePath)
	}

	// Verify data in users table by checking for returned errors
	hasData := uu[0].GetByEmail()
	if hasData == nil {
		log.Info().Msgf("[%s] env already has seeded data", currentEnv)
		return nil
	}

	for i, u := range uu {
		pw, err := u.HashPasswordUsingBcrypt()
		if err != nil {
			log.Error().Err(err).Msg("bcrypt_error")
			return err_rest.NewInternalServerError(("something went wrong"))
		}

		uu[i].Password = pw
		uu[i].Save()
	}

	log.Info().Msgf("successfully seeded data in [%s] env", currentEnv)
	return nil
}
*/

/*
func CloseConnections() {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	if err := appInstance.ShutdownWithContext(ctx); err != nil {
		log.Err(err).Msg("failed to gracefully shutdown the server")
	}

	cancel()
	log.Info().Msg("app instance has shutdown")

	rdbmsInstance.Disconnect()
	inMemoryDbInstance.Disconnect()
}
*/

/*
func loadUsersFromJsonFile(filePath string) ([]user.User, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var users []user.User
	if err := json.Unmarshal(byteValue, &users); err != nil {
		return nil, err
	}

	return users, nil
}
*/

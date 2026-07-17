package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/DarrelA/starter-go-postgresql/internal/domain/apperror"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/repository"
	"github.com/DarrelA/starter-go-postgresql/internal/infrastructure/auth"
	"github.com/rs/zerolog/log"
)

const (
	envBasePath                = "/root/deployment/build"
	errMsgUnableToLoadJSONFile = "unable to load [%s]"
)

type PostgresSeedRepository struct {
	env         string
	envBasePath string
}

func NewSeedRepository(env string) *PostgresSeedRepository {
	return &PostgresSeedRepository{env: env, envBasePath: envBasePath}
}

func (sr PostgresSeedRepository) Seed(ctx context.Context, ur repository.UserRepository) error {
	currentEnv := sr.env
	switch currentEnv {
	case "dev":
		return saveMultipleUsers(ctx, currentEnv, sr.envBasePath, ur)
	case "test":
		return saveMultipleUsers(ctx, currentEnv, sr.envBasePath, ur)
	default:
		log.Info().Msgf("[%s] env will NOT be seeded with data", currentEnv)
		return nil
	}
}

func saveMultipleUsers(
	ctx context.Context,
	currentEnv string,
	envBasePath string,
	ur repository.UserRepository,
) error {
	userJsonFilePath := "/seed.user." + currentEnv + ".json"
	uu, err := loadUsersFromJsonFile(envBasePath + "/json" + userJsonFilePath)
	if err != nil {
		return fmt.Errorf(errMsgUnableToLoadJSONFile+": %w", userJsonFilePath, err)
	}
	if len(uu) == 0 {
		return fmt.Errorf("seed file %s contains no users", userJsonFilePath)
	}

	seeded := 0
	passwordService := auth.NewPasswordService()
	for _, u := range uu {
		_, err := ur.GetByEmail(ctx, u.Email)
		if err == nil {
			continue
		}
		if !errors.Is(err, apperror.ErrUserNotFound) {
			return fmt.Errorf("check seed user %s: %w", u.Email, err)
		}

		hashedPassword, err := passwordService.Hash(u.Password)
		if err != nil {
			return fmt.Errorf("hash password for seed user %s: %w", u.Email, err)
		}
		u.Password = hashedPassword
		if err := ur.Save(ctx, u); err != nil {
			// Another process may have inserted the same seed user after the
			// existence check. Treat that race as an idempotent success.
			if errors.Is(err, apperror.ErrEmailConflict) {
				continue
			}
			return fmt.Errorf("save seed user %s: %w", u.Email, err)
		}
		seeded++
	}

	log.Info().Int("users_added", seeded).Msgf("successfully reconciled seed data in [%s] env", currentEnv)
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

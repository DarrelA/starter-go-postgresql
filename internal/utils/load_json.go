package utils

import (
	"encoding/json"
	"io"
	"os"

	user "github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
)

func LoadUsersFromJsonFile(filePath string) ([]user.User, error) {
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

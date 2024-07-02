package utils

import (
	"encoding/json"
	"io"
	"os"

	"github.com/DarrelA/starter-go-postgresql/internal/domains/users"
)

func LoadUsersFromJsonFile(filePath string) ([]users.User, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var users []users.User
	if err := json.Unmarshal(byteValue, &users); err != nil {
		return nil, err
	}

	return users, nil
}

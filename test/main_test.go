package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/DarrelA/starter-go-postgresql/app"
	"github.com/DarrelA/starter-go-postgresql/db"
	"github.com/gofiber/fiber/v2"
)

var (
	rdbmsInstance       db.RDBMS
	inMemoryDbInstance  db.InMemoryDB
	appInstance         *fiber.App
	authServiceInstance *fiber.App
)

func TestMain(m *testing.M) {
	// Change to the project root directory from `/test` directory
	err := os.Chdir("../")
	if err != nil {
		fmt.Println("error changing directory:", err)
		os.Exit(1)
	}

	rdbmsInstance, inMemoryDbInstance = app.CreateDBConnections()
	appInstance, authServiceInstance = app.ConfigureAppInstance()
	go app.StartServer()

	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestMainFunc(t *testing.T) {
	logFilePath := "./deployments/logs/app.log"

	t.Run("Test if app.log is writable", func(t *testing.T) {
		if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
			t.Fatalf("expected log file to be created, but it does not exist")
		}

		file, err := os.OpenFile(logFilePath, os.O_WRONLY, 0644)
		if err != nil {
			t.Fatalf("expected log file to be writable, but it is not: %v", err)
		}

		file.Close()
	})
}

func TestFiberHandler(t *testing.T) {
	// @TODO: Build and load `baseURL` in `configs\app`
	protocol := "http://"
	baseURL := "localhost:8080/auth/api"

	input := map[string]string{
		"first_name": "a",
		"last_name":  "1",
		"email":      "e@e.com",
		"password":   "1",
	}

	body, _ := json.Marshal(input)

	req, err := http.NewRequest(http.MethodPost, protocol+baseURL+"/register", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
	}

	var responseBody map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&responseBody)
}

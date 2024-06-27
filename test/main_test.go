package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/DarrelA/starter-go-postgresql/app"
	"github.com/DarrelA/starter-go-postgresql/configs"
	"github.com/DarrelA/starter-go-postgresql/db"
	"github.com/DarrelA/starter-go-postgresql/internal/domains/users"
	"github.com/gofiber/fiber/v2"
)

var (
	rdbmsInstance       db.RDBMS
	inMemoryDbInstance  db.InMemoryDB
	appInstance         *fiber.App
	authServiceInstance *fiber.App
)

var baseURL = configs.BaseURLs.AuthService

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

	t.Run("test if app.log is writable", func(t *testing.T) {
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

func TestAuthService(t *testing.T) {
	var registerInputs = []users.RegisterInput{
		{FirstName: "kang yao", LastName: "tan", Email: "tky@e.com", Password: "I_love_c_sharp"},
		{FirstName: "jie wei", LastName: "low", Email: "ljw@e.com", Password: "i_love_java"},
		{FirstName: "bing hong", LastName: "tan", Email: "tbh@e.com", Password: "i_heArt_VB"},
		{FirstName: "jason", LastName: "the consultant", Email: "j@e.com", Password: "I&asked(a)question^at!the*town-hall:_Why@is9the6air~conditioning%so+cold"},
	}

	// Using a single HTTP client for all requests
	client := &http.Client{}

	for _, input := range registerInputs {
		t.Run(fmt.Sprintf("test inserting [%s %s] into rdbms", input.FirstName, input.LastName), func(t *testing.T) {
			body, err := json.Marshal(input)
			if err != nil {
				t.Fatalf("failed to marshal json: %v", err)
			}

			req, err := http.NewRequest(http.MethodPost, baseURL+"/register", bytes.NewBuffer(body))
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			req.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("failed to send request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
			}

			// Decode the response body into a temporary map
			var responseMap map[string]json.RawMessage
			if err := json.NewDecoder(resp.Body).Decode(&responseMap); err != nil {
				t.Errorf("failed to decode response body into map: %v", err)
			}

			// Extract the "user" field and decode it into UserResponse
			var responseBody users.UserResponse
			if err := json.Unmarshal(responseMap["user"], &responseBody); err != nil {
				t.Errorf("failed to decode user field: %v", err)
			}

			if responseBody.UUID == nil {
				t.Errorf("expected UUID to be created for [%s %s], but it is empty", input.FirstName, input.LastName)
			}
			if responseBody.FirstName != input.FirstName {
				t.Errorf("expected FirstName to be [%s], but got [%s]", input.FirstName, responseBody.FirstName)
			}
			if responseBody.LastName != input.LastName {
				t.Errorf("expected LastName to be [%s], but got [%s]", input.LastName, responseBody.LastName)
			}
			if responseBody.Email != input.Email {
				t.Errorf("expected Email to be [%s], but got [%s]", input.Email, responseBody.Email)
			}
		})
	}
}

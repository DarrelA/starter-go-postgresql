package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/DarrelA/starter-go-postgresql/app"
	"github.com/DarrelA/starter-go-postgresql/configs"
	"github.com/DarrelA/starter-go-postgresql/db"
	"github.com/DarrelA/starter-go-postgresql/internal/domains/users"
	"github.com/DarrelA/starter-go-postgresql/internal/utils/err_rest"
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
			cwd, err := os.Getwd()
			if err != nil {
				t.Fatalf("cannot get cwd: %s", err)
			}

			t.Logf("@cwd: %s", cwd)
			t.Fatal("expected log file to be created, but it does not exist")
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

		// all empty fields
		{FirstName: "", LastName: "", Email: "", Password: ""},
		// empty password only
		{FirstName: "John", LastName: "Doe", Email: "John_Doe@gmail.com", Password: ""},
		// empty email only
		{FirstName: "Jane", LastName: "Smith", Email: "", Password: "Password1!"},
		// empty last name only
		{FirstName: "Alice", LastName: "", Email: "Alice@yahoo.com", Password: "Password1!"},
		// empty first name only
		{FirstName: "", LastName: "Brown", Email: "Brown@outlook.com", Password: "Password1!"},
		// valid input
		{FirstName: "Emily", LastName: "Clark", Email: "Emily_Clark@gmail.com", Password: "Password1!"},
		// email is already taken
		{FirstName: "emily", LastName: "clark", Email: "emily_clark@gmail.com", Password: "Password1!"},
		// invalid email format
		{FirstName: "Oliver", LastName: "Jones", Email: "Oliver_Jones", Password: "Password1!"},
		// invalid password (too short)
		{FirstName: "Michael", LastName: "Taylor", Email: "Michael_Taylor@yahoo.com", Password: "Pass1!"},
		// invalid password (no special character)
		{FirstName: "Emma", LastName: "Davis", Email: "Emma_Davis@outlook.com", Password: "Password1"},
		// invalid password (no number)
		{FirstName: "William", LastName: "Martinez", Email: "William_Martinez@gmail.com", Password: "Password!"},
		// invalid first name (less than 2 characters)
		{FirstName: "A", LastName: "Anderson", Email: "A_Anderson@yahoo.com", Password: "Password1!"},
		// invalid last name (less than 2 characters)
		{FirstName: "Sophia", LastName: "B", Email: "Sophia_B@outlook.com", Password: "Password1!"},
		// invalid first name (more than 50 characters)
		{FirstName: "Taumatawhakatangihangakoauauotamateaturipukakapikimaungahoronukupokaiwhenu" + "12345678901234567890", LastName: "Harris", Email: "LongFirstName@gmail.com", Password: "Password1!"},
		// invalid last name (more than 50 characters)
		{FirstName: "James", LastName: "Taumatawhakatangihangakoauauotamateaturipukakapikimaungahoronukupokaiwhenu" + "12345678901234567890", Email: "LongLastName@gmail.com", Password: "Password1!"},
		// invalid first name (contains non-alphabetic characters)
		{FirstName: "John123", LastName: "Walker", Email: "John123_Walker@gmail.com", Password: "Password1!"},
		// invalid last name (contains non-alphabetic characters)
		{FirstName: "Grace", LastName: "Miller456", Email: "Grace_Miller456@yahoo.com", Password: "Password1!"},
	}

	// Using a single HTTP client for all requests
	client := &http.Client{}

	t.Log("Testing /register endpoint")
	for _, user := range registerInputs {
		t.Run(fmt.Sprintf("test inserting [%s %s] into rdbms", user.FirstName, user.LastName), func(t *testing.T) {
			body, err := json.Marshal(user)
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

			// Decode the response body into a temporary map
			var responseMap map[string]json.RawMessage
			if err := json.NewDecoder(resp.Body).Decode(&responseMap); err != nil {
				t.Errorf("failed to decode response body into map: %v", err)
			}

			t.Log("Test for StatusOK")
			if resp.StatusCode == http.StatusOK {
				// Extract the "user" field and decode it into UserResponse
				var responseBody users.UserResponse
				if err := json.Unmarshal(responseMap["user"], &responseBody); err != nil {
					t.Errorf("failed to decode field: %v", err)
				}

				if responseBody.UUID == nil {
					t.Errorf("expected UUID to be created for [%s %s], but it is empty", user.FirstName, user.LastName)
				}
				if responseBody.FirstName != user.FirstName {
					t.Errorf("expected FirstName to be [%s], but got [%s]", user.FirstName, responseBody.FirstName)
				}
				if responseBody.LastName != user.LastName {
					t.Errorf("expected LastName to be [%s], but got [%s]", user.LastName, responseBody.LastName)
				}
				if responseBody.Email != user.Email {
					t.Errorf("expected Email to be [%s], but got [%s]", user.Email, responseBody.Email)
				}
			}

			t.Log("Test for StatusBadRequest")
			if resp.StatusCode == http.StatusBadRequest {
				var responseBody err_rest.RestErr
				if err := json.Unmarshal(responseMap["error"], &responseBody); err != nil {
					t.Errorf("failed to decode field: %v", err)
				}

				if !strings.Contains(responseBody.Message, "validation error") &&
					!strings.Contains(responseBody.Message, "email is already taken") {
					t.Errorf("expected message to contain \"validation error\" or \"email is already taken\", but got [%s]", responseBody.Message)
				}
			}

			// Expecting no StatusInternalServerError since client should receive only StatusBadRequest or StatusOK
			t.Log("Test for StatusInternalServerError")
			if resp.StatusCode == http.StatusInternalServerError {
				var responseBody err_rest.RestErr
				if err := json.Unmarshal(responseMap["error"], &responseBody); err != nil {
					t.Errorf("failed to decode field: %v", err)
				}

				t.Errorf("expected ZERO StatusInternalServerError, but got [%s]", responseBody.Message)
			}
		})
	}
}

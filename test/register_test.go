package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/DarrelA/starter-go-postgresql/configs"
	user "github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	"github.com/DarrelA/starter-go-postgresql/internal/utils/err_rest"
	data_test "github.com/DarrelA/starter-go-postgresql/test/data"
)

/*
var (
	rdbmsInstance       db.RDBMS
	inMemoryDbInstance  db.InMemoryDB
	appInstance         *fiber.App
	authServiceInstance *fiber.App
)
*/

func TestMain(m *testing.M) {
	// envs_utils.CreateAppLog()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		// rdbmsInstance, inMemoryDbInstance = app.CreateDBConnections()
		// app.SeedDatabase()
		// appInstance, authServiceInstance = app.ConfigureAppInstance()
		// go app.StartServer()
	}()

	wg.Wait()
	exitVal := m.Run()
	// envs_utils.GetLogFile().Close()
	os.Exit(exitVal)
}

func TestRegisterEndpoint(t *testing.T) {
	// Using a single HTTP client for all requests
	client := &http.Client{}
	baseURL := configs.BaseURLs.AuthService
	endpoint := "/register"

	for _, newUser := range data_test.RegisterInputs {
		t.Run(fmt.Sprintf("test case for [%s]: ", newUser.TestName), func(t *testing.T) {
			// Extract the RegisterInput for the HTTP request
			// userRegisterInput := user.RegisterInput
			// body, err := json.Marshal(userRegisterInput)
			body, err := json.Marshal(newUser)
			if err != nil {
				t.Fatalf("failed to marshal json: %v", err)
			}

			req, err := http.NewRequest(http.MethodPost, baseURL+endpoint, bytes.NewBuffer(body))
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			req.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("failed to send request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != newUser.ExpectedStatusCode {
				t.Fatalf("expected Status Code to be [%d], but got [%d]", newUser.ExpectedStatusCode, resp.StatusCode)
			}

			// Decode the response body into a temporary map
			var responseMap map[string]json.RawMessage
			if err := json.NewDecoder(resp.Body).Decode(&responseMap); err != nil {
				t.Errorf("failed to decode response body into map: %v", err)
			}

			switch resp.StatusCode {
			case http.StatusOK:
				// Extract the "user" field and decode it into UserResponse
				var responseBody user.UserResponse
				if err := json.Unmarshal(responseMap["user"], &responseBody); err != nil {
					t.Errorf("failed to decode field: %v", err)
				}

				if responseBody.UUID == nil {
					t.Errorf("expected UUID to be created for [%s %s], but it is empty", newUser.FirstName, newUser.LastName)
				}
				if responseBody.FirstName != strings.TrimSpace(strings.ToLower(newUser.FirstName)) {
					t.Errorf("expected FirstName to be [%s], but got [%s]", newUser.FirstName, responseBody.FirstName)
				}
				if responseBody.LastName != strings.TrimSpace(strings.ToLower(newUser.LastName)) {
					t.Errorf("expected LastName to be [%s], but got [%s]", newUser.LastName, responseBody.LastName)
				}
				if responseBody.Email != strings.TrimSpace(strings.ToLower(newUser.Email)) {
					t.Errorf("expected Email to be [%s], but got [%s]", newUser.Email, responseBody.Email)
				}

			case http.StatusBadRequest:
				var responseBody err_rest.RestErr
				if err := json.Unmarshal(responseMap["error"], &responseBody); err != nil {
					t.Errorf("failed to decode field: %v", err)
				}

				if !strings.Contains(responseBody.Message, "validation error") &&
					!strings.Contains(responseBody.Message, "email is already taken") {
					t.Errorf("expected message to contain \"validation error\" or \"email is already taken\", but got [%s]", responseBody.Message)
				}

			// Expecting no StatusInternalServerError since client should receive only StatusBadRequest or StatusOK
			case http.StatusInternalServerError:
				var responseBody err_rest.RestErr
				if err := json.Unmarshal(responseMap["error"], &responseBody); err != nil {
					t.Fatalf("failed to decode field: %v", err)
				}

				t.Errorf("expected ZERO StatusInternalServerError, but got [%s]", responseBody.Message)

			default:
				t.Errorf("unexpected error [Status Code - %d]: [%s]", resp.StatusCode, responseMap)
			}
		})
	}
}

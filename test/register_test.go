package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/DarrelA/starter-go-postgresql/configs"
	"github.com/DarrelA/starter-go-postgresql/internal/domains/users"
	"github.com/DarrelA/starter-go-postgresql/internal/utils/err_rest"
	data_test "github.com/DarrelA/starter-go-postgresql/test/data"
)

func TestRegisterEndpoint(t *testing.T) {
	// Using a single HTTP client for all requests
	client := &http.Client{}
	baseURL := configs.BaseURLs.AuthService

	for _, user := range data_test.RegisterInputs {
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

			t.Log("test for StatusOK")
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

			t.Log("test for StatusBadRequest")
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
			t.Log("test for StatusInternalServerError")
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

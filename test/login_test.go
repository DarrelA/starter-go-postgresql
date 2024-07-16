package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/DarrelA/starter-go-postgresql/internal/utils/err_rest"
	data_test "github.com/DarrelA/starter-go-postgresql/test/data"
)

func TestLoginEndpoint(t *testing.T) {
	// Using a single HTTP client for all requests
	client := &http.Client{}
	baseURL := "@TODO: Fix tests"
	endpoint := "/login"

	for _, user := range data_test.LoginInputs {
		t.Run(fmt.Sprintf("test case for [%s]", user.TestName), func(t *testing.T) {
			body, err := json.Marshal(user)
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

			if resp.StatusCode != user.ExpectedStatusCode {
				t.Fatalf("expected Status Code to be [%d], but got [%d]", resp.StatusCode, user.ExpectedStatusCode)
			}

			// Decode the response body into a temporary map
			var responseMap map[string]json.RawMessage
			if err := json.NewDecoder(resp.Body).Decode(&responseMap); err != nil {
				t.Errorf("failed to decode response body into map: %v", err)
			}

			switch resp.StatusCode {
			case http.StatusOK:
				var accessToken string
				if err := json.Unmarshal(responseMap["access_token"], &accessToken); err != nil {
					t.Fatalf("failed to unmarshal access token: %v", err)
				}

				if accessToken == "" {
					t.Errorf("expected AccessToken to be created for [%s], but it is empty", user.Email)
				}

			case http.StatusBadRequest:
				var responseBody err_rest.RestErr
				if err := json.Unmarshal(responseMap["error"], &responseBody); err != nil {
					t.Errorf("failed to decode field: %v", err)
				}

				if !strings.Contains(responseBody.Message, "validation error") && !strings.Contains(responseBody.Message, "invalid credentials") &&
					!strings.Contains(responseBody.Message, "the account has not been registered") {
					t.Errorf("expected message to contain \"validation error\" or \"invalid credentials\" or \"the account has not been registered\", but got [%s]", responseBody.Message)
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

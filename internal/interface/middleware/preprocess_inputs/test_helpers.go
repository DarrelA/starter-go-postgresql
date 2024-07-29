// coverage:ignore file
// Test file
package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v2"
)

// registerHandler handles the register route
func registerHandler(c *fiber.Ctx) error {
	payload := c.Locals("register_payload")
	if payload == nil {
		return c.Status(fiber.StatusBadRequest).SendString("No register payload found")
	}
	return c.JSON(payload)
}

// loginHandler handles the login route
func loginHandler(c *fiber.Ctx) error {
	payload := c.Locals("login_payload")
	if payload == nil {
		return c.Status(fiber.StatusBadRequest).SendString("No login payload found")
	}
	return c.JSON(payload)
}

// createRequest creates a new test request based on the given test case
func createRequest(t *testing.T, test testCase) *http.Request {
	if test.payload != nil {
		payloadBytes, err := json.Marshal(test.payload)
		if err != nil {
			t.Errorf("Failed to marshal payload: %v", err)
		}
		req := httptest.NewRequest("POST", test.url, bytes.NewReader(payloadBytes))
		req.Header.Set("Content-Type", "application/json")
		return req
	}
	return httptest.NewRequest("POST", test.url, nil)
}

// createDummyRequest creates a new dummy request based on the test case
func createDummyRequest(t *testing.T, test testCase) *http.Request {
	if test.invalidJSON {
		return httptest.NewRequest("POST", "/dummy", bytes.NewReader([]byte(`{invalid json`)))
	}
	payloadBytes, err := json.Marshal(test.payload)
	if err != nil {
		t.Errorf("failed to marshal payload: %v", err)
	}
	req := httptest.NewRequest("POST", "/dummy", bytes.NewReader(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	return req
}

// assertResponse asserts the response of the dummy request
func assertResponse(t *testing.T, resp *http.Response, test testCase) {
	if resp.StatusCode != test.expectedStatus {
		t.Errorf("expected status %d but got %d", test.expectedStatus, resp.StatusCode)
	}

	if test.invalidJSON {
		assertInvalidJSONResponse(t, resp, test.expectedError)
	} else {
		assertValidJSONResponse(t, resp, test.expectedPayload)
	}
}

// assertInvalidJSONResponse asserts the response for an invalid JSON request
func assertInvalidJSONResponse(t *testing.T, resp *http.Response, expectedError string) {
	var respBody map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		t.Errorf("failed to decode response body: %v", err)
	}

	errorMsg, ok := respBody["error"].(map[string]interface{})["message"].(string)
	if !ok {
		t.Errorf("failed to get error message from response body")
	}

	if errorMsg != expectedError {
		t.Errorf("expected error message '%s' but got '%s'", expectedError, errorMsg)
	}
}

// assertValidJSONResponse asserts the response for a valid JSON request
func assertValidJSONResponse(t *testing.T, resp *http.Response, expectedPayload interface{}) {
	result := reflect.New(reflect.TypeOf(expectedPayload)).Interface()
	err := json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	if !reflect.DeepEqual(expectedPayload, reflect.ValueOf(result).Elem().Interface()) {
		t.Errorf("expected payload %+v but got %+v", expectedPayload, reflect.ValueOf(result).Elem().Interface())
	}
}

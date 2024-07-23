// coverage:ignore file
package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	"github.com/gofiber/fiber/v2"
)

// setupFiberTestApp sets up the test application with necessary routes and middleware
func setupFiberTestApp() *fiber.App {
	app := fiber.New()

	// A middleware that mocks setting baseURLsConfig in Locals
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("baseURLsConfig", &entity.BaseURLsConfig{
			AuthServicePathName: authServicePathName,
		})
		return c.Next()
	})

	app.Use(PreProcessInputs)

	app.Post(authServicePathName+"/register", registerHandler)
	app.Post(authServicePathName+"/login", loginHandler)

	return app
}

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

// createPreProcessInputsTestRequest creates a new test request based on the given test case
func createPreProcessInputsTestRequest(t *testing.T, test testCase) *http.Request {
	if test.payload != nil {
		payloadBytes, err := json.Marshal(test.payload)
		if err != nil {
			t.Fatalf("Failed to marshal payload: %v", err)
		}
		req := httptest.NewRequest("POST", test.url, bytes.NewReader(payloadBytes))
		req.Header.Set("Content-Type", "application/json")
		return req
	}
	return httptest.NewRequest("POST", test.url, nil)
}

// performAppTest performs the test request on the given app
func performAppTest(t *testing.T, app *fiber.App, req *http.Request) *http.Response {
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("An error occurred: %v", err)
	}
	return resp
}

// assertEmail checks if the email in the response body matches the expected email
func assertEmail(t *testing.T, body []byte, expectedEmail string) {
	var responseBody map[string]interface{}
	err := json.Unmarshal(body, &responseBody)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	email, ok := responseBody["email"].(string)
	if !ok {
		t.Fatalf("Failed to get email from response body")
	}
	if email != expectedEmail {
		t.Fatalf("Expected email '%s' but got '%s'", expectedEmail, email)
	}
}

// createDummyRoute creates a dummy route for testing
func createDummyRoute(app *fiber.App, test testCase) {
	app.Post("/dummy", func(c *fiber.Ctx) error {
		if test.invalidJSON {
			var invalidPayload interface{}
			if err := parseAndSanitize(c, invalidPayload); err != nil {
				return err
			}
		} else {
			payload := reflect.New(reflect.TypeOf(test.payload).Elem()).Interface()
			if err := parseAndSanitize(c, payload); err != nil {
				return err
			}
			return c.JSON(payload)
		}
		return nil
	})
}

// createDummyRequest creates a new dummy request based on the test case
func createDummyRequest(t *testing.T, test testCase) *http.Request {
	if test.invalidJSON {
		return httptest.NewRequest("POST", "/dummy", bytes.NewReader([]byte(`{invalid json`)))
	}
	payloadBytes, err := json.Marshal(test.payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}
	req := httptest.NewRequest("POST", "/dummy", bytes.NewReader(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	return req
}

// createParseAndSanitizeTestRequest performs a request to the dummy route
func createParseAndSanitizeTestRequest(t *testing.T, app *fiber.App, test testCase) *http.Response {
	req := createDummyRequest(t, test)
	return performAppTest(t, app, req)
}

// assertResponse asserts the response of the dummy request
func assertResponse(t *testing.T, resp *http.Response, test testCase) {
	if resp.StatusCode != test.expectedStatus {
		t.Fatalf("expected status %d but got %d", test.expectedStatus, resp.StatusCode)
	}

	if test.invalidJSON {
		assertInvalidJSONResponse(t, resp, test.expectedError)
	} else {
		assertValidJSONResponse(t, resp, test.expectedPayload)
	}
}

// assertInvalidJSONResponse asserts the response for an invalid JSON request
func assertInvalidJSONResponse(t *testing.T, resp *http.Response, expectedError string) {
	var responseBody map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	errorMessage, ok := responseBody["error"].(map[string]interface{})["message"].(string)
	if !ok {
		t.Fatalf("failed to get error message from response body")
	}

	if errorMessage != expectedError {
		t.Fatalf("expected error message '%s' but got '%s'", expectedError, errorMessage)
	}
}

// assertValidJSONResponse asserts the response for a valid JSON request
func assertValidJSONResponse(t *testing.T, resp *http.Response, expectedPayload interface{}) {
	result := reflect.New(reflect.TypeOf(expectedPayload)).Interface()
	err := json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !reflect.DeepEqual(expectedPayload, reflect.ValueOf(result).Elem().Interface()) {
		t.Fatalf("expected payload %+v but got %+v", expectedPayload, reflect.ValueOf(result).Elem().Interface())
	}
}

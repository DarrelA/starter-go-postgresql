package middleware

import (
	"bytes"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestLoggerMW(t *testing.T) {
	tests := []struct {
		name                  string
		env                   any
		requestID             any
		correlationID         any
		expectedSubString     string
		expectedRequestID     string
		expectedCorrelationID string
	}{
		{
			name: "currentEnvProd", env: "prod", requestID: "test-request-id", correlationID: "test-correlation-id", expectedSubString: "request is completed in [prod] env",
			expectedRequestID: "test-request-id", expectedCorrelationID: "test-correlation-id",
		},
		{
			name: "currentEnvDev", env: "dev", requestID: "test-request-id", correlationID: "test-correlation-id", expectedSubString: "request is completed in [dev] env",
			expectedRequestID: "", expectedCorrelationID: "",
		},
		{
			name: "currentEnvDefault", env: "test", requestID: "test-request-id", correlationID: "test-correlation-id", expectedSubString: "request is completed in [test] env",
			expectedRequestID: "test-request-id", expectedCorrelationID: "test-correlation-id",
		},
		{
			name: "invalidRequestID", env: "test", requestID: 0, correlationID: "test-correlation-id",
			expectedSubString: "status 400: type error: request_id",
			expectedRequestID: "", expectedCorrelationID: "",
		},
		{
			name: "invalidCorrelationID", env: "test", requestID: "test-request-id", correlationID: 0,
			expectedSubString: "status 400: type error: correlation_id",
			expectedRequestID: "", expectedCorrelationID: "",
		},
		{
			name: "invalidCurrentEnv", env: 0, requestID: "test-request-id", correlationID: "test-correlation-id",
			expectedSubString: "status 400: type error: currentEnv",
			expectedRequestID: "", expectedCorrelationID: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, buf := setupAppAndBuffer(tt.env, tt.requestID, tt.correlationID)

			// Create a test request
			req := httptest.NewRequest("GET", "/", nil)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("LoggerMW test failed: %v", err)
			}

			// Ensure response status code is 200 OK
			if resp.StatusCode != fiber.StatusOK {
				t.Fatalf("Expected status '%d' but got '%d'", fiber.StatusOK, resp.StatusCode)
			}

			checkLogOutput(t, buf, tt.expectedSubString, tt.expectedRequestID, tt.expectedCorrelationID)
		})
	}
}

func TestGetIPVersion(t *testing.T) {
	tests := []struct {
		ip            string
		expectedValue string
	}{
		{"", "unknown"},
		{"5.59.32.2", "IPv4"},
		{"5bf7:a43f:402d:e01d:7b5d:071f:5068:62b7", "IPv6"},
	}

	for _, tt := range tests {
		value := getIPVersion(tt.ip)
		if value != tt.expectedValue {
			t.Fatalf("Expected value '%s' but got '%s'", tt.expectedValue, value)
		}
	}
}

func setupAppAndBuffer(env any, requestID any, correlationID any) (
	*fiber.App, *bytes.Buffer) {
	app := fiber.New()

	// Set up a buffer to capture log output
	var buf bytes.Buffer
	log.Logger = log.Output(&buf)
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Add required context to the app via middleware
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("request_id", requestID)
		c.Locals("correlation_id", correlationID)
		c.Locals("env", env)
		return c.Next()
	})
	app.Use(LoggerMW)

	// Define a test route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	return app, &buf
}

func checkLogOutput(
	t *testing.T, buf *bytes.Buffer, expectedSubString string,
	expectedRequestID string, expectedCorrelationID string,
) {
	logOutput := buf.String()
	if !strings.Contains(logOutput, expectedSubString) {
		t.Errorf("Expected '%s' not found in log output: %s", expectedSubString, logOutput)
	}
	if expectedRequestID != "" && !strings.Contains(logOutput, expectedRequestID) {
		t.Errorf("Expected request_id '%s' not found in log output: %s", expectedRequestID, logOutput)
	}
	if expectedCorrelationID != "" && !strings.Contains(logOutput, expectedCorrelationID) {
		t.Errorf("Expected correlation_id '%s' not found in log output: %s", expectedCorrelationID, logOutput)
	}
	// Reset the logger to its default output
	log.Logger = log.Output(os.Stdout)
}

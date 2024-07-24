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
			expectedSubString: "status 400: type error: requestID",
			expectedRequestID: "", expectedCorrelationID: "",
		},
		{
			name: "invalidCorrelationID", env: "test", requestID: "test-request-id", correlationID: 0,
			expectedSubString: "status 400: type error: correlationID",
			expectedRequestID: "", expectedCorrelationID: "",
		},
		{
			name: "invalidCurrentEnv", env: 0, requestID: "test-request-id", correlationID: "test-correlation-id",
			expectedSubString: "status 400: type error: currentEnv",
			expectedRequestID: "", expectedCorrelationID: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			app := fiber.New()

			// Set up a buffer to capture log output
			var buf bytes.Buffer
			log.Logger = log.Output(&buf)
			zerolog.SetGlobalLevel(zerolog.InfoLevel)

			// Add required context to the app via middleware
			app.Use(func(c *fiber.Ctx) error {
				c.Locals("requestID", test.requestID)
				c.Locals("correlationID", test.correlationID)
				c.Locals("env", test.env)
				return c.Next()
			})

			// Define a test route
			app.Get("/", LoggerMW, func(c *fiber.Ctx) error {
				return c.SendString("Hello, World!")
			})

			// Create a test request
			req := httptest.NewRequest("GET", "/", nil)
			resp, err := app.Test(req)
			if err != nil {
				t.Errorf("LoggerMW test failed: %v", err)
			}

			// Ensure response status code is 200 OK
			if resp.StatusCode != fiber.StatusOK {
				t.Errorf("Expected status '%d' but got '%d'", fiber.StatusOK, resp.StatusCode)
			}

			logOutput := buf.String()
			if !strings.Contains(logOutput, test.expectedSubString) {
				t.Errorf("Expected '%s' not found in log output: %s", test.expectedSubString, logOutput)
			}
			if test.expectedRequestID != "" && !strings.Contains(logOutput, test.expectedRequestID) {
				t.Errorf("Expected requestID '%s' not found in log output: %s", test.expectedRequestID, logOutput)
			}
			if test.expectedCorrelationID != "" && !strings.Contains(logOutput, test.expectedCorrelationID) {
				t.Errorf("Expected correlationID '%s' not found in log output: %s", test.expectedCorrelationID, logOutput)
			}
			// Reset the logger to its default output
			log.Logger = log.Output(os.Stdout)
		})
	}

	t.Run("Test getIPVersion", func(t *testing.T) {
		tests := []struct {
			ip            string
			expectedValue string
		}{
			{"", "unknown"},
			{"5.59.32.2", "IPv4"},
			{"5bf7:a43f:402d:e01d:7b5d:071f:5068:62b7", "IPv6"},
		}

		for _, test := range tests {
			value := getIPVersion(test.ip)
			if value != test.expectedValue {
				t.Errorf("Expected value '%s' but got '%s'", test.expectedValue, value)
			}
		}
	})
}

package middleware

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	"github.com/gofiber/fiber/v2"
)

// TestPreProcessInputs tests the PreProcessInputs middleware
func TestPreProcessInputs(t *testing.T) {
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

	for _, test := range preProcessInputsTests {
		t.Run(test.name, func(t *testing.T) {
			req := createRequest(t, test)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("An error occurred: %v", err)
			}

			defer resp.Body.Close()

			var responseBody map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			if err != nil {
				t.Fatalf("Failed to decode response body: %v", err)
			}

			if test.expectedErrorMsg != "" {
				errorMsg, ok := responseBody["error"].(map[string]interface{})["message"].(string)
				if !ok || !strings.Contains(errorMsg, test.expectedErrorMsg) {
					t.Errorf("Expected error message to contain %q, got %q", test.expectedErrorMsg, errorMsg)
				}
			} else {
				if resp.StatusCode != fiber.StatusOK {
					t.Errorf("Expected status code %d, got %d.", fiber.StatusOK, resp.StatusCode)
				}

				email, ok := responseBody["email"].(string)
				if !ok {
					t.Fatalf("Failed to get email from response body")
				}
				if email != test.expectedEmail {
					t.Fatalf("Expected email '%s' but got '%s'", test.expectedEmail, email)
				}
			}
		})
	}

	t.Run("Test normalizePath", func(t *testing.T) {
		for _, test := range normalizePathTests {
			t.Run(test.name, func(t *testing.T) {
				normalizePath := normalizePath(test.path)
				if normalizePath != test.expectedValue {
					t.Fatalf("Expected value '%s' but got '%s'", test.expectedValue, normalizePath)
				}
			})
		}
	})
}

// TestParseAndSanitize tests the parseAndSanitize function
func TestParseAndSanitize(t *testing.T) {
	for _, test := range parseAndSanitizeTests {
		t.Run(test.name, func(t *testing.T) {
			app := fiber.New()
			app.Post("/dummy", func(c *fiber.Ctx) error {
				if test.invalidJSON {
					var invalidPayload interface{}
					if err := parseAndSanitize(c, invalidPayload); err != nil {
						return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "error": err})
					}
				} else {
					payload := reflect.New(reflect.TypeOf(test.payload).Elem()).Interface()
					if err := parseAndSanitize(c, payload); err != nil {
						return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "error": err})
					}
					return c.JSON(payload)
				}
				return nil
			})

			req := createDummyRequest(t, test)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("An error occurred: %v", err)
			}

			assertResponse(t, resp, test)
		})
	}

	t.Run("Test sanitizeHelper", func(t *testing.T) {
		for _, test := range sanitizeHelperTests {
			t.Run(test.name, func(t *testing.T) {
				sanitizeHelper(&test.registerPayload)
				sanitizeHelper(&test.loginPayload)
				if test.registerPayload != test.expectedSanitizedRP {
					t.Fatalf("Expected value '%s' but got '%s'", test.expectedSanitizedRP, test.registerPayload)
				}
				if test.loginPayload != test.expectedSanitizedLP {
					t.Fatalf("Expected value '%s' but got '%s'", test.expectedSanitizedLP, test.loginPayload)
				}
			})
		}
	})
}

func TestValidateStruct(t *testing.T) {
	/*
	   1.	Payload Definition: The `payload` in your test struct is defined as `*TestPayload`.
	   2.	Passing by Reference: When initializing `payload` for each test case, use `&TestPayload{}` to create a pointer.
	   3.	Function Call: Pass `test.payload` directly to `validateStruct` since `test.payload` is already a pointer.
	*/
	for _, test := range testValidateStructTests {
		t.Run(test.name, func(t *testing.T) {
			err := validateStruct(test.payload)
			if test.expectedErrMsg == "" {
				// Expect no error
				if err != nil {
					t.Fatalf("Expected no error but got '%s'", err.Error())
				}
			} else {
				if err == nil {
					t.Fatalf("Expected an error but got nil")
				}
				if !strings.Contains(err.Error(), test.expectedErrMsg) {
					t.Fatalf("Expected error message to contain '%s' but got '%s'", test.expectedErrMsg, err.Error())
				}
			}
		})
	}

	t.Run("Test camelToSnakeCase", func(t *testing.T) {
		for _, test := range camelToSnakeCaseTests {
			value := camelToSnakeCase(test.str)
			if value != test.expectedValue {
				t.Fatalf("Expected value '%s' but got '%s'", test.expectedValue, value)
			}
		}
	})
}

func TestValidatePassword(t *testing.T) {
	for _, test := range validatePasswordTests {
		// Create a mock `FieldLevel` to pass into the `validatePassword` function
		mockFieldLevel := &mockFieldLevel{
			field: test.password,
		}
		value := validatePassword(mockFieldLevel)
		if value != test.expectedValue {
			t.Fatalf("Expected value '%v' but got '%v' for password '%s'", test.expectedValue, value, test.password)
		}
	}
}

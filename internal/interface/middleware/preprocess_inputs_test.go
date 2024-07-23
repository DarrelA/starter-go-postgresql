package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	dto "github.com/DarrelA/starter-go-postgresql/internal/application/repository/dto"
	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	"github.com/gofiber/fiber/v2"
)

const authServicePathName = "/auth/api/v1/users"

func TestPreProcessInputs(t *testing.T) {
	app := fiber.New()

	// Middleware to mock setting baseURLsConfig in Locals
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("baseURLsConfig", &entity.BaseURLsConfig{
			AuthServicePathName: authServicePathName,
		})
		return c.Next()
	})

	// PreProcessInputs middleware
	app.Use(PreProcessInputs)

	// Dummy register route
	app.Post(authServicePathName+"/register", func(c *fiber.Ctx) error {
		payload := c.Locals("register_payload")
		if payload == nil {
			return c.Status(fiber.StatusBadRequest).SendString("No register payload found")
		}
		return c.JSON(payload)
	})

	// Dummy login route
	app.Post(authServicePathName+"/login", func(c *fiber.Ctx) error {
		payload := c.Locals("login_payload")
		if payload == nil {
			return c.Status(fiber.StatusBadRequest).SendString("No login payload found")
		}
		return c.JSON(payload)
	})

	tests := []struct {
		name             string
		url              string
		payload          interface{}
		expectedEmail    string
		expectedErrorMsg string
	}{
		{
			name: "Valid register endpoint",
			url:  authServicePathName + "/register",
			payload: dto.RegisterInput{
				FirstName: "Jie",
				LastName:  "Wei",
				Email:     "JieWei@gmail.com",
				Password:  "P@ssword1",
			},
			expectedEmail: "jiewei@gmail.com",
		},
		{
			name: "Valid login endpoint",
			url:  authServicePathName + "/login",
			payload: dto.LoginInput{
				Email:    "JieWei@gmail.com",
				Password: "P@ssword1",
			},
			expectedEmail: "jiewei@gmail.com",
		},
		{
			name:             "Invalid endpoint",
			url:              "/auth/invalid",
			expectedErrorMsg: "invalid endpoint: /auth/invalid",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var req *http.Request

			if test.payload != nil {
				payloadBytes, err := json.Marshal(test.payload)
				if err != nil {
					t.Fatalf("Failed to marshal payload: %v", err)
				}
				req = httptest.NewRequest("POST", test.url, bytes.NewReader(payloadBytes))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = httptest.NewRequest("POST", test.url, nil)
			}

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("An error occurred: %v", err)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %v", err)
			}

			// Check for expected error message in the response body
			if test.expectedErrorMsg != "" {
				if !strings.Contains(string(body), test.expectedErrorMsg) {
					t.Errorf("Expected error message to contain %q, got %q", test.expectedErrorMsg, string(body))
				}
			}

			// Check for the expected status code
			if resp.StatusCode != fiber.StatusOK && test.expectedErrorMsg == "" {
				t.Errorf("Expected status code %d, got %d.", fiber.StatusOK, resp.StatusCode)
			}

			// Decode the response body
			var responseBody map[string]interface{}
			err = json.Unmarshal(body, &responseBody)
			if err != nil {
				t.Fatalf("Failed to decode response body: %v", err)
			}

			if test.expectedErrorMsg == "" {
				// Extract and verify the email from the response body
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
		tests := []struct {
			name          string
			path          string
			expectedValue string
		}{
			{
				name:          "NoExtraSlashes",
				path:          "localhost:8080/auth/api/v1/users/register",
				expectedValue: "localhost:8080/auth/api/v1/users/register",
			},
			{
				name:          "ExtraTrailingSlashes",
				path:          "localhost:8080/auth/api/v1/users/register///////",
				expectedValue: "localhost:8080/auth/api/v1/users/register",
			},
			{
				name:          "ExtraTrailingBackslashes",
				path:          "localhost:8080/auth/api/v1/users/register\\\\\\\\",
				expectedValue: "localhost:8080/auth/api/v1/users/register",
			},
			{
				name:          "MixedSlashesAndBackslashes",
				path:          "localhost:8080/auth/api/v1/users/register//\\/\\//\\///\\\\",
				expectedValue: "localhost:8080/auth/api/v1/users/register",
			},
			{
				name:          "BackslashesInPath",
				path:          "localhost:8080\\auth\\api\\v1\\users\\register",
				expectedValue: "localhost:8080/auth/api/v1/users/register",
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				normalizePath := normalizePath(test.path)
				if normalizePath != test.expectedValue {
					t.Fatalf("Expected value '%s' but got '%s'", test.expectedValue, normalizePath)
				}
			})
		}
	})
}

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

type testCase struct {
	name            string
	payload         interface{}
	expectedPayload interface{}
	invalidJSON     bool
	expectedStatus  int
	expectedError   string
}

func performRequest(t *testing.T, app *fiber.App, test testCase) *http.Response {
	var req *http.Request
	if test.invalidJSON {
		// Create a new HTTP request with invalid JSON
		req = httptest.NewRequest("POST", "/dummy", bytes.NewReader([]byte(`{invalid json`)))
	} else {
		// Create a new HTTP request with the payload
		payloadBytes, err := json.Marshal(test.payload)
		if err != nil {
			t.Fatalf("failed to marshal payload: %v", err)
		}
		req = httptest.NewRequest("POST", "/dummy", bytes.NewReader(payloadBytes))
	}
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	resp, err := app.Test(req, 3)
	if err != nil {
		t.Fatalf("failed to perform request: %v", err)
	}
	return resp
}

func assertResponse(t *testing.T, resp *http.Response, test testCase) {
	// Assert the response status code
	if resp.StatusCode != test.expectedStatus {
		t.Fatalf("expected status %d but got %d", test.expectedStatus, resp.StatusCode)
	}

	if test.invalidJSON {
		// Decode the response body for invalid JSON case
		var responseBody map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&responseBody)
		if err != nil {
			t.Fatalf("failed to decode response body: %v", err)
		}

		// Extract the error message
		errorMessage, ok := responseBody["error"].(map[string]interface{})["message"].(string)
		if !ok {
			t.Fatalf("failed to get error message from response body")
		}

		// Assert the response body contains the expected error message
		if errorMessage != test.expectedError {
			t.Fatalf("expected error message '%s' but got '%s'", test.expectedError, errorMessage)
		}
	} else {
		// Decode the response payload into the expected payload type
		result := reflect.New(reflect.TypeOf(test.expectedPayload)).Interface()
		err := json.NewDecoder(resp.Body).Decode(result)
		if err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		// Assert the response payload matches the expected payload
		if !reflect.DeepEqual(test.expectedPayload, reflect.ValueOf(result).Elem().Interface()) {
			t.Fatalf("expected payload %+v but got %+v", test.expectedPayload, reflect.ValueOf(result).Elem().Interface())
		}
	}
}

func TestParseAndSanitize(t *testing.T) {
	tests := []testCase{
		{
			name:            "Valid RegisterInput",
			payload:         &dto.RegisterInput{FirstName: "", LastName: "", Email: "", Password: ""},
			expectedPayload: dto.RegisterInput{FirstName: "", LastName: "", Email: "", Password: ""},
			expectedStatus:  fiber.StatusOK,
		},
		{
			name: "Sanitized RegisterInput",
			payload: &dto.RegisterInput{
				FirstName: "    Jie    ", LastName: "    Wei    ", Email: "    JieWei@gmail.com", Password: "P@ssword1    ",
			},
			expectedPayload: dto.RegisterInput{
				FirstName: "jie", LastName: "wei", Email: "jiewei@gmail.com", Password: "P@ssword1",
			},
			expectedStatus: fiber.StatusOK,
		},
		{
			name:            "Valid LoginInput",
			payload:         &dto.LoginInput{Email: "", Password: ""},
			expectedPayload: dto.LoginInput{Email: "", Password: ""},
			expectedStatus:  fiber.StatusOK,
		},
		{
			name:            "Sanitized LoginInput",
			payload:         &dto.LoginInput{Email: "    JieWei@gmail.com    ", Password: "    P@ssword1    "},
			expectedPayload: dto.LoginInput{Email: "jiewei@gmail.com", Password: "P@ssword1"},
			expectedStatus:  fiber.StatusOK,
		},
		{
			name:           "Invalid JSON",
			payload:        nil, // Payload is not used in this case
			invalidJSON:    true,
			expectedStatus: fiber.StatusUnprocessableEntity,
			expectedError:  errInvalidJSON,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			app := fiber.New()
			createDummyRoute(app, test)
			resp := performRequest(t, app, test)
			assertResponse(t, resp, test)
		})
	}

	t.Run("Test sanitizeHelper", func(t *testing.T) {
		tests := []struct {
			name                string
			registerPayload     dto.RegisterInput
			loginPayload        dto.LoginInput
			expectedSanitizedRP dto.RegisterInput
			expectedSanitizedLP dto.LoginInput
		}{
			{
				name:                "EmptyPayloads",
				registerPayload:     dto.RegisterInput{FirstName: "", LastName: "", Email: "", Password: ""},
				loginPayload:        dto.LoginInput{Email: "", Password: ""},
				expectedSanitizedRP: dto.RegisterInput{FirstName: "", LastName: "", Email: "", Password: ""},
				expectedSanitizedLP: dto.LoginInput{Email: "", Password: ""},
			},
			{
				name: "ValidPayloads",
				registerPayload: dto.RegisterInput{
					FirstName: "Jie", LastName: "Wei", Email: "JieWei@gmail.com", Password: "P@ssword1",
				},
				loginPayload: dto.LoginInput{
					Email: "JieWei@gmail.com", Password: "P@ssword1",
				},
				expectedSanitizedRP: dto.RegisterInput{
					FirstName: "jie", LastName: "wei", Email: "jiewei@gmail.com", Password: "P@ssword1",
				},
				expectedSanitizedLP: dto.LoginInput{
					Email: "jiewei@gmail.com", Password: "P@ssword1",
				},
			},
			{
				name: "PayloadsWithTrailingSpaces",
				registerPayload: dto.RegisterInput{
					FirstName: "Jie    ", LastName: "Wei    ", Email: "JieWei@gmail.com    ", Password: "P@ssword1    ",
				},
				loginPayload: dto.LoginInput{
					Email: "JieWei@gmail.com    ", Password: "P@ssword1    ",
				},
				expectedSanitizedRP: dto.RegisterInput{
					FirstName: "jie", LastName: "wei", Email: "jiewei@gmail.com", Password: "P@ssword1",
				},
				expectedSanitizedLP: dto.LoginInput{
					Email: "jiewei@gmail.com", Password: "P@ssword1",
				},
			},
			{
				name: "PayloadsWithLeadingSpaces",
				registerPayload: dto.RegisterInput{
					FirstName: "    Jie", LastName: "    Wei", Email: "    JieWei@gmail.com", Password: "    P@ssword1",
				},
				loginPayload: dto.LoginInput{
					Email: "    JieWei@gmail.com", Password: "    P@ssword1",
				},
				expectedSanitizedRP: dto.RegisterInput{
					FirstName: "jie", LastName: "wei", Email: "jiewei@gmail.com", Password: "P@ssword1",
				},
				expectedSanitizedLP: dto.LoginInput{
					Email: "jiewei@gmail.com", Password: "P@ssword1",
				},
			},
			{
				name: "PayloadsWithUppercaseLetters",
				registerPayload: dto.RegisterInput{
					FirstName: "JIE", LastName: "WEI", Email: "JIEWEI@GMAIL.COM", Password: "P@ssword1",
				},
				loginPayload: dto.LoginInput{
					Email: "JIEWEI@GMAIL.COM", Password: "P@ssword1",
				},
				expectedSanitizedRP: dto.RegisterInput{
					FirstName: "jie", LastName: "wei", Email: "jiewei@gmail.com", Password: "P@ssword1",
				},
				expectedSanitizedLP: dto.LoginInput{
					Email: "jiewei@gmail.com", Password: "P@ssword1",
				},
			},
		}

		for _, test := range tests {
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
	type TestPayload struct {
		RequiredField string `validate:"required"`
		MinField      string `validate:"min=5"`
		MaxField      string `validate:"max=10"`
		AlphaField    string `validate:"alpha"`
		EmailField    string `validate:"email"`
		PasswdField   string `validate:"passwd"`
	}

	tests := []struct {
		name           string
		payload        *TestPayload
		expectedErrMsg string
	}{
		{
			name: "Required field missing",
			payload: &TestPayload{
				MinField:    "12345",
				MaxField:    "1234567890",
				AlphaField:  "abcde",
				EmailField:  "test@example.com",
				PasswdField: "Password1!",
			},
			expectedErrMsg: fmt.Sprintf("the field [%s] should %s", "required_field", requiredVM),
		},
		{
			name: "Min field too short",
			payload: &TestPayload{
				RequiredField: "required",
				MinField:      "123",
				MaxField:      "1234567890",
				AlphaField:    "abcde",
				EmailField:    "test@example.com",
				PasswdField:   "Password1!",
			},
			expectedErrMsg: fmt.Sprintf("the field [%s] should "+minVM, "min_field", "5"),
		},
		{
			name: "Max field too long",
			payload: &TestPayload{
				RequiredField: "required",
				MinField:      "12345",
				MaxField:      "12345678901",
				AlphaField:    "abcde",
				EmailField:    "test@example.com",
				PasswdField:   "Password1!",
			},
			expectedErrMsg: fmt.Sprintf("the field [%s] should "+maxVM, "max_field", "10"),
		},
		{
			name: "Alpha field with non-alphabetic characters",
			payload: &TestPayload{
				RequiredField: "required",
				MinField:      "12345",
				MaxField:      "1234567890",
				AlphaField:    "abc123",
				EmailField:    "test@example.com",
				PasswdField:   "Password1!",
			},
			expectedErrMsg: fmt.Sprintf("the field [%s] should %s", "alpha_field", alphaVM),
		},
		{name: "Invalid email",
			payload: &TestPayload{
				RequiredField: "required",
				MinField:      "12345",
				MaxField:      "1234567890",
				AlphaField:    "abcde",
				EmailField:    "invalid-email",
				PasswdField:   "Password1!",
			},
			expectedErrMsg: fmt.Sprintf("the field [%s] should %s", "email_field", emailVM),
		},
		{
			name: "Invalid password",
			payload: &TestPayload{
				RequiredField: "required",
				MinField:      "12345",
				MaxField:      "1234567890",
				AlphaField:    "abcde",
				EmailField:    "test@example.com",
				PasswdField:   "password",
			},
			expectedErrMsg: fmt.Sprintf("the field [%s] should %s", "passwd_field", passwdVM),
		},
		{
			name: "Invalid email and invalid password",
			payload: &TestPayload{
				RequiredField: "required",
				MinField:      "12345",
				MaxField:      "1234567890",
				AlphaField:    "abcde",
				EmailField:    "invalid-email",
				PasswdField:   "password",
			},
			expectedErrMsg: fmt.Sprintf("the field [%s] should %s\n", "email_field", emailVM) +
				fmt.Sprintf("the field [%s] should %s", "passwd_field", passwdVM),
		},
		{
			name: "Valid payload",
			payload: &TestPayload{
				RequiredField: "required",
				MinField:      "12345",
				MaxField:      "12345",
				AlphaField:    "abcde",
				EmailField:    "test@example.com",
				PasswdField:   "Password1!",
			},
		},
	}

	/*
	   1.	Payload Definition: The `payload` in your test struct is defined as `*TestPayload`.
	   2.	Passing by Reference: When initializing `payload` for each test case, use `&TestPayload{}` to create a pointer.
	   3.	Function Call: Pass `test.payload` directly to `validateStruct` since `test.payload` is already a pointer.
	*/
	for _, test := range tests {
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
		tests := []struct {
			str           string
			expectedValue string
		}{
			{"pleaseChangeFromCamelToSnakeCase", "please_change_from_camel_to_snake_case"},
			{"thisIsATest", "this_is_a_test"},
			{"convertThisString", "convert_this_string"},
			{"CamelCase", "camel_case"},
			{"lowercase", "lowercase"},
			{"!aA@bB#cC$%^&*()_-+=", "!a_a@b_b#c_c$%^&*()_-+="},
		}

		for _, test := range tests {
			value := camelToSnakeCase(test.str)
			if value != test.expectedValue {
				t.Fatalf("Expected value '%s' but got '%s'", test.expectedValue, value)
			}
		}
	})
}

// MockFieldLevel implements validator.FieldLevel for testing purposes
type MockFieldLevel struct {
	field string
}

func (m *MockFieldLevel) Top() reflect.Value {
	return reflect.Value{}
}

func (m *MockFieldLevel) Parent() reflect.Value {
	return reflect.Value{}
}

func (m *MockFieldLevel) Field() reflect.Value {
	return reflect.ValueOf(m.field)
}

func (m *MockFieldLevel) FieldName() string {
	return ""
}

func (m *MockFieldLevel) StructFieldName() string {
	return ""
}

func (m *MockFieldLevel) Param() string {
	return ""
}

func (m *MockFieldLevel) GetTag() string {
	return ""
}

func (m *MockFieldLevel) ExtractType(field reflect.Value) (reflect.Value, reflect.Kind, bool) {
	return field, field.Kind(), false
}

func (m *MockFieldLevel) GetStructFieldOK() (reflect.Value, reflect.Kind, bool) {
	return reflect.Value{}, reflect.Invalid, false
}

func (m *MockFieldLevel) GetStructFieldOKAdvanced(val reflect.Value, namespace string) (reflect.Value, reflect.Kind, bool) {
	return reflect.Value{}, reflect.Invalid, false
}

func (m *MockFieldLevel) GetStructFieldOK2() (reflect.Value, reflect.Kind, bool, bool) {
	return reflect.Value{}, reflect.Invalid, false, false
}

func (m *MockFieldLevel) GetStructFieldOKAdvanced2(val reflect.Value, namespace string) (reflect.Value, reflect.Kind, bool, bool) {
	return reflect.Value{}, reflect.Invalid, false, false
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		password      string
		expectedValue bool
	}{
		{"Password1!", true},
		{"password1", false},
		{"PASSWORD1", false},
		{"Password", false},
		{"P@ssword", false},
		{"P@ssword1", true},
	}

	for _, test := range tests {
		// Create a mock FieldLevel to pass into the validatePassword function
		mockFieldLevel := &MockFieldLevel{
			field: test.password,
		}
		value := validatePassword(mockFieldLevel)
		if value != test.expectedValue {
			t.Fatalf("Expected value '%v' but got '%v' for password '%s'", test.expectedValue, value, test.password)
		}
	}
}

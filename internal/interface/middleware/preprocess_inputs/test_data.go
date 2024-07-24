// coverage:ignore file
package middleware

import (
	"fmt"

	dto "github.com/DarrelA/starter-go-postgresql/internal/application/repository/dto"
	"github.com/gofiber/fiber/v2"
)

const authServicePathName = "/auth/api/v1/users"

type testCase struct {
	name             string
	url              string
	payload          interface{}
	expectedEmail    string
	expectedErrorMsg string
	expectedStatus   int
	expectedError    string
	expectedPayload  interface{}
	invalidJSON      bool
}

var preProcessInputsTests = []testCase{
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
	{
		name:             "Failed to parse and sanitize register payload",
		url:              authServicePathName + "/register",
		payload:          nil, // This would simulate invalid JSON or payload
		expectedErrorMsg: "invalid json body",
	},
	{
		name:             "Failed to parse and sanitize login payload",
		url:              authServicePathName + "/login",
		payload:          nil, // This would simulate invalid JSON or payload
		expectedErrorMsg: "invalid json body",
	},
	{
		name: "Failed to validate register payload",
		url:  authServicePathName + "/register",
		payload: dto.RegisterInput{
			FirstName: "", LastName: "", Email: "invalidEmail", Password: "",
		},
		expectedErrorMsg: fmt.Sprintf("the field [%s] should %s\n", "first_name", requiredVM) +
			fmt.Sprintf("the field [%s] should %s\n", "last_name", requiredVM) +
			fmt.Sprintf("the field [%s] should %s\n", "email", emailVM) +
			fmt.Sprintf("the field [%s] should %s", "password", requiredVM),
	},
	{
		name: "Failed to validate login payload",
		url:  authServicePathName + "/login",
		payload: dto.LoginInput{
			Email: "invalidEmail", Password: "",
		},
		expectedErrorMsg: fmt.Sprintf("the field [%s] should %s\n", "email", emailVM) +
			fmt.Sprintf("the field [%s] should %s", "password", requiredVM),
	},
}

var normalizePathTests = []struct {
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

var parseAndSanitizeTests = []testCase{
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

var sanitizeHelperTests = []struct {
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

type testPayload struct {
	RequiredField string `validate:"required"`
	MinField      string `validate:"min=5"`
	MaxField      string `validate:"max=10"`
	AlphaField    string `validate:"alpha"`
	EmailField    string `validate:"email"`
	PasswdField   string `validate:"passwd"`
}

var testValidateStructTests = []struct {
	name           string
	payload        *testPayload
	expectedErrMsg string
}{
	{
		name: "Required field missing",
		payload: &testPayload{
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
		payload: &testPayload{
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
		payload: &testPayload{
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
		payload: &testPayload{
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
		payload: &testPayload{
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
		payload: &testPayload{
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
		payload: &testPayload{
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
		payload: &testPayload{
			RequiredField: "required",
			MinField:      "12345",
			MaxField:      "12345",
			AlphaField:    "abcde",
			EmailField:    "test@example.com",
			PasswdField:   "Password1!",
		},
	},
}

var camelToSnakeCaseTests = []struct {
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

var validatePasswordTests = []struct {
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

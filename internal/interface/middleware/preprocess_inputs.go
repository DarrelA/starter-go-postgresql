package middlewares

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"unicode"

	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	dto "github.com/DarrelA/starter-go-postgresql/internal/interface/transport/dto"
	"github.com/DarrelA/starter-go-postgresql/internal/utils/err_rest"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func PreProcessInputs(c *fiber.Ctx) error {
	baseURLsConfig, ok := c.Locals("baseURLsConfig").(*entity.BaseURLsConfig)
	if !ok {
		log.Panic().Msg("invalid_baseURLsConfig")
	}

	authServicePathName := baseURLsConfig.AuthServicePathName
	endpoint := normalizePath(c.Path())

	switch endpoint {
	case authServicePathName + "/register":
		var payload dto.RegisterInput
		if err := parseAndSanitize(c, &payload); err != nil {
			return err
		}

		if err := validateStruct(&payload); err != nil {
			return err
		}

		c.Locals("register_payload", payload)

	case authServicePathName + "/login":
		var payload dto.LoginInput
		if err := parseAndSanitize(c, &payload); err != nil {
			return err
		}

		if err := validateStruct(&payload); err != nil {
			return err
		}

		c.Locals("login_payload", payload)

	default:
		log.Error().Msg("invalid endpoint for PreProcessInputs middleware")
		err := err_rest.NewBadRequestError("invalid endpoint: " + endpoint)
		return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "error": err})
	}

	return c.Next()
}

// normalizePath removes trailing slashes
func normalizePath(path string) string {
	path = strings.ReplaceAll(path, "\\", "/")
	normalized := strings.TrimRight(path, "/")
	return normalized
}

func parseAndSanitize(c *fiber.Ctx, payload interface{}) error {
	if err := c.BodyParser(payload); err != nil {
		err := err_rest.NewBadRequestError("invalid json body")
		return c.Status(err.Status).JSON(fiber.Map{"status": "fail", "error": err})
	}

	sanitizeHelper(payload)
	return nil
}

/*
Reflection allows you to inspect and modify the values of variables at runtime.
`reflect.ValueOf(payload)` gets the reflection value of the payload.
`.Elem()` gets the underlying value that the pointer points to.
`v.NumField()` returns the number of fields in the struct.
`v.Field(i)` accesses each field by its index.
`field.Kind()` returns the kind of the field (e.g., string, int).
*/

// sanitizeHelper sanitizes the input struct by trimming spaces and converting strings to lowercase
func sanitizeHelper(payload interface{}) {
	v := reflect.ValueOf(payload).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.String {
			field.SetString(strings.TrimSpace(field.String()))
			if !strings.Contains(strings.ToLower(v.Type().Field(i).Name), "password") {
				field.SetString(strings.ToLower(field.String()))
			}
		}
	}
}

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterValidation("passwd", validatePassword)
}

var validationMessages = map[string]string{
	"required": "not be empty",
	"min":      "be at least %s characters long",
	"max":      "be at most %s characters long",
	"alpha":    "contain only alphabetic characters",
	"alphanum": "contain only alphanumeric characters",
	"email":    "be a valid email address",
	"passwd":   "contain at least one number, one uppercase letter, one lowercase letter, and one special character",
}

/*
For each validation error, the corresponding user-friendly message is fetched from the validationMessages map.
If a tag does not have a predefined message, a default message "be valid" is used.
The message is constructed using the field name and the user-friendly message.

Pass the payload by reference
*/
func validateStruct[T any](payload *T) *err_rest.RestErr {
	var validationErrors []string
	err := validate.Struct(payload)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			tag := err.Tag()
			messageTemplate, exists := validationMessages[tag]

			if !exists {
				messageTemplate = "be valid"
			}

			fieldName := camelToSnakeCase(err.Field())
			message := ""

			if strings.Contains(messageTemplate, "%s") {
				message = fmt.Sprintf("the field [%s] should "+messageTemplate, fieldName, err.Param())
			} else {
				message = fmt.Sprintf("the field [%s] should %s", fieldName, messageTemplate)
			}

			validationErrors = append(validationErrors, message)
		}

		fullMessage := "validation error: " + strings.Join(validationErrors, "\n")
		return err_rest.NewBadRequestError(fullMessage)
	}

	return nil
}

func camelToSnakeCase(str string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	var hasNumber, hasUpper, hasLower, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasNumber && hasUpper && hasLower && hasSpecial
}

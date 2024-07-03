package middlewares

import (
	"reflect"
	"strings"

	"github.com/DarrelA/starter-go-postgresql/configs"
	"github.com/DarrelA/starter-go-postgresql/internal/domains/users"
	"github.com/DarrelA/starter-go-postgresql/internal/utils/err_rest"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func PreProcessInputs(c *fiber.Ctx) error {
	authServicePathName := configs.BaseURLs.AuthServicePathName
	endpoint := normalizePath(c.Path())

	switch endpoint {
	case authServicePathName + "/register":
		var payload users.RegisterInput
		if err := parseAndSanitize(c, &payload); err != nil {
			return err
		}

		c.Locals("register_payload", payload)

	case authServicePathName + "/login":
		var payload users.LoginInput
		if err := parseAndSanitize(c, &payload); err != nil {
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

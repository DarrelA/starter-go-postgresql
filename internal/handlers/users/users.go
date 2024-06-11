package users

import (
	"github.com/DarrelA/starter-go-postgresql/internal/domains/users"
	"github.com/DarrelA/starter-go-postgresql/internal/services"
	"github.com/DarrelA/starter-go-postgresql/internal/utils/errors"
	"github.com/gofiber/fiber/v2"
)

func Register(c *fiber.Ctx) error {
	var payload users.UserRequest

	if err := c.BodyParser(&payload); err != nil {
		err := errors.NewBadRequestError("invalid json body")
		return c.Status(err.Status).JSON(err)
	}

	result, saveErr := services.CreateUser(payload)
	if saveErr != nil {
		return c.Status(saveErr.Status).JSON(saveErr)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

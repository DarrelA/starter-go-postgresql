package users

import (
	"github.com/DarrelA/starter-go-postgresql/internal/domains/users"
	"github.com/DarrelA/starter-go-postgresql/internal/services"
	"github.com/DarrelA/starter-go-postgresql/internal/utils/errors"
	"github.com/gofiber/fiber/v2"
)

func Register(c *fiber.Ctx) error {
	var payload users.RegisterInput

	if err := c.BodyParser(&payload); err != nil {
		err := errors.NewBadRequestError("invalid json body")
		return c.Status(err.Status).JSON(err)
	}

	result, err := services.CreateUser(payload)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

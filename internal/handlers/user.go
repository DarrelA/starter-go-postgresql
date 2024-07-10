package handlers

import (
	"github.com/DarrelA/starter-go-postgresql/internal/domain/users"
	"github.com/gofiber/fiber/v2"
)

func GetMe(c *fiber.Ctx) error {
	user_record := c.Locals("user_record").(*users.UserRecord)
	return c.Status(fiber.StatusOK).
		JSON(fiber.Map{"status": "success", "data": fiber.Map{"user_record": user_record}})
}

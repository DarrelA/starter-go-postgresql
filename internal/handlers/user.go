package handlers

import (
	dto "github.com/DarrelA/starter-go-postgresql/internal/interface/transport/dto"
	"github.com/gofiber/fiber/v2"
)

func GetMe(c *fiber.Ctx) error {
	user_record := c.Locals("user_record").(*dto.UserRecord)
	return c.Status(fiber.StatusOK).
		JSON(fiber.Map{"status": "success", "data": fiber.Map{"user_record": user_record}})
}

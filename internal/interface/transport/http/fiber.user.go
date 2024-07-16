package http

import (
	"github.com/DarrelA/starter-go-postgresql/internal/application/service"
	dto "github.com/DarrelA/starter-go-postgresql/internal/interface/transport/dto"
	"github.com/gofiber/fiber/v2"
)

// The `UserService` struct serves both as the business logic implementation and the HTTP adapter.
type UserService struct{}

func NewUserService() service.UserService {
	return &UserService{}
}

func (us *UserService) GetUserRecord(c *fiber.Ctx) error {
	user_record := c.Locals("user_record").(*dto.UserRecord)
	return c.Status(fiber.StatusOK).
		JSON(fiber.Map{"status": "success", "data": fiber.Map{"user_record": user_record}})
}

package http

import (
	dto "github.com/DarrelA/starter-go-postgresql/internal/application/repository/dto"
	"github.com/DarrelA/starter-go-postgresql/internal/application/service"
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

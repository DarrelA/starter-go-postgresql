// coverage:ignore file
// Testing with integration test
package http

import (
	dto "github.com/DarrelA/starter-go-postgresql/internal/application/dto"
	"github.com/DarrelA/starter-go-postgresql/internal/application/usecase"
	"github.com/gofiber/fiber/v2"
)

type UserUseCase struct{}

func NewUserUseCase() usecase.UserUseCase {
	return &UserUseCase{}
}

func (uuc *UserUseCase) GetUserRecord(c *fiber.Ctx) error {
	userRecord := c.Locals("userRecord").(*dto.UserRecord)
	return c.Status(fiber.StatusOK).
		JSON(fiber.Map{"status": "success", "data": fiber.Map{"userRecord": userRecord}})
}

package usecase

import "github.com/gofiber/fiber/v2"

type UserUseCase interface {
	GetUserRecord(c *fiber.Ctx) error
}

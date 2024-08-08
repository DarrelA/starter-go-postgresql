package usecase

import "github.com/gofiber/fiber/v2"

type AuthUseCase interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	RefreshAccessToken(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
}

type OAuth2UseCase interface {
	Login(c *fiber.Ctx) error
	Callback(c *fiber.Ctx) error
}

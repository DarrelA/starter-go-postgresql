package app

import (
	"github.com/DarrelA/starter-go-postgresql/internal/handlers"
	"github.com/DarrelA/starter-go-postgresql/internal/middlewares"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func authServiceRouter(router *fiber.App) {
	v1 := router.Group("/api")

	v1.Post("/register", handlers.Register)
	v1.Post("/login", handlers.Login)
	v1.Get("/logout", middlewares.Deserializer, handlers.Logout)
	v1.Get("/refresh", handlers.RefreshAccessToken)
	v1.Get("/users/me", middlewares.Deserializer, handlers.GetMe)

	v1.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
	})

	v1.All("*", func(c *fiber.Ctx) error {
		path := c.Path()
		log.Error().Msg("Invalid Path: " + path)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "404 - Not Found",
		})
	})
}

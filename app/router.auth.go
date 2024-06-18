package app

import (
	"github.com/DarrelA/starter-go-postgresql/internal/handlers"
	"github.com/DarrelA/starter-go-postgresql/internal/middlewares"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func authServiceRouter() {
	authService.Route("/auth", func(router fiber.Router) {
		router.Post("/register", handlers.Register)
		router.Post("/login", handlers.Login)
		router.Get("/logout", handlers.Logout)
		router.Get("/refresh", handlers.RefreshAccessToken)
	})

	authService.Get("/users/me", middlewares.Deserializer, handlers.GetMe)

	authService.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
	})

	authService.All("*", func(c *fiber.Ctx) error {
		path := c.Path()
		log.Error().Msg("Invalid Path: " + path)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "404 - Not Found",
		})
	})
}

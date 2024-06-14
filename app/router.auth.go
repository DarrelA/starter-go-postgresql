package app

import (
	"log"

	"github.com/DarrelA/starter-go-postgresql/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

func authServiceRouter() {
	authService.Route("/auth", func(router fiber.Router) {
		router.Post("/register", handlers.Register)
		router.Post("/login", handlers.Login)
		router.Get("/refresh", handlers.RefreshAccessToken)
	})

	authService.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
	})

	authService.All("*", func(c *fiber.Ctx) error {
		path := c.Path()
		log.Printf("Invalid Path: %v", path)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "404 - Not Found",
		})
	})
}
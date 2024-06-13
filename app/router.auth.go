package app

import (
	"github.com/DarrelA/starter-go-postgresql/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

func authServiceRouter() {
	authService.Route("/auth", func(router fiber.Router) {
		router.Post("/register", handlers.Register)
		router.Post("/login", handlers.Login)
		router.Get("/refresh", handlers.RefreshAccessToken)
	})
}

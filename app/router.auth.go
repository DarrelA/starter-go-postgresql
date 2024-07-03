package app

import (
	"github.com/DarrelA/starter-go-postgresql/internal/handlers"
	mw "github.com/DarrelA/starter-go-postgresql/internal/middlewares"
	"github.com/DarrelA/starter-go-postgresql/internal/utils/err_rest"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func authServiceRouter(router *fiber.App) {
	v1 := router.Group("/api")

	v1.Post("/register", mw.PreProcessInputs, handlers.Register)
	v1.Post("/login", mw.PreProcessInputs, handlers.Login)
	v1.Get("/logout", mw.Deserializer, handlers.Logout)
	v1.Get("/refresh", handlers.RefreshAccessToken)
	v1.Get("/users/me", mw.Deserializer, handlers.GetMe)

	v1.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
	})

	v1.All("*", func(c *fiber.Ctx) error {
		path := c.Path()
		err := err_rest.NewBadRequestError("Invalid Path: " + path)
		log.Error().Err(err).Msg("")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "404 - Not Found",
		})
	})
}

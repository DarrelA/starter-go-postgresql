package app

import (
	"github.com/DarrelA/starter-go-postgresql/configs"
	"github.com/DarrelA/starter-go-postgresql/internal/handlers/users"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func mapUrls() {
	app.Use(cors.New(cors.Config{
		AllowOrigins:     configs.CORSSettings.AllowedOrigins,
		AllowMethods:     "GET,POST",
		AllowHeaders:     "Content-Type",
		ExposeHeaders:    "Content-Length",
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60,
	}))

	app.Post("/api/register", users.Register)
	app.Post("/api/login", users.Login)
}

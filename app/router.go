package app

import (
	"time"

	"github.com/DarrelA/starter-go-postgresql/configs"
	"github.com/DarrelA/starter-go-postgresql/internal/handlers/users"
	"github.com/gin-contrib/cors"
)

func mapUrls() {
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{configs.CORSSettings.AllowedOrigins},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == configs.CORSSettings.AllowedOrigins
		},
		MaxAge: 12 * time.Hour,
	}))

	router.POST("/api/register", users.Register)
	router.POST("/api/login", users.Login)
	router.GET("/api/user", users.Get)
	router.GET("/api/logout", users.Logout)
}

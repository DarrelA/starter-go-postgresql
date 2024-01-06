package app

import (
	"github.com/DarrelA/starter-go-postgresql/internal/handlers/users"
)

func mapUrls() {
	router.POST("/api/register", users.Register)
	router.POST("/api/login", users.Login)
	router.GET("/api/user", users.Get)
}

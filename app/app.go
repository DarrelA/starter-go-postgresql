package app

import (
	"github.com/DarrelA/starter-go-postgresql/configs"
	db "github.com/DarrelA/starter-go-postgresql/db/pgdb"
	"github.com/gofiber/fiber/v2"
)

var (
	app = fiber.New()
)

func StartApp() {
	db.ConnectPostgresDatabase()
	mapUrls()
	app.Listen(":" + configs.Port)
}

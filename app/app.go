package app

import (
	"github.com/DarrelA/starter-go-postgresql/configs"
	"github.com/DarrelA/starter-go-postgresql/db/pgdb"
	"github.com/gofiber/fiber/v2"
)

var (
	app = fiber.New()
)

func StartApp() {
	pgdb.ConnectDatabase()
	mapUrls()
	app.Listen(":" + configs.Port)
}

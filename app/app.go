package app

import (
	"github.com/DarrelA/starter-go-postgresql/configs"
	"github.com/DarrelA/starter-go-postgresql/db/pgdb"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func StartApp() {
	pgdb.ConnectDatabase()
	mapUrls()
	router.Run(":" + configs.Port)
}

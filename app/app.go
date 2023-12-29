package app

import (
	pgdb "github.com/DarrelA/starter-go-postgresql/server/db"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func StartApp() {
	pgdb.ConnectDatabase()
	router.Run(":4040")
}

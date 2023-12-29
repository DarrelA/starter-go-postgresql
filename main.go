package main

import (
	"log"

	"github.com/DarrelA/starter-go-postgresql/server/app"
)

func main() {
	app.StartApp()
	log.Println("hello")
}

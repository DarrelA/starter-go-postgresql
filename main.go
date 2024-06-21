package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/DarrelA/starter-go-postgresql/app"
	"github.com/DarrelA/starter-go-postgresql/configs"
	"github.com/DarrelA/starter-go-postgresql/db"
	pgdb "github.com/DarrelA/starter-go-postgresql/db/pgdb"
	redisDb "github.com/DarrelA/starter-go-postgresql/db/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Debug().Msg("Starting main function")

	// Create a logger instance with output to a file
	logFile, err := os.Create("/app/logs/app.log")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create log file")
	}

	// Configure logger to write to both file and console
	log.Logger = zerolog.
		New(zerolog.
			MultiLevelWriter(zerolog.ConsoleWriter{Out: os.Stderr}, logFile)).
		With().
		Caller().
		Timestamp().
		Logger()

	// Configure Instances
	rdbmsInstance := pgdb.NewDB()
	inMemoryDbInstance := redisDb.NewDB()
	appInstance := app.StartApp(rdbmsInstance, inMemoryDbInstance)

	go startServer(appInstance)
	waitForShutdown(appInstance, rdbmsInstance, inMemoryDbInstance)

	logFile.Close()
	os.Exit(0)
}

func startServer(appInstance *fiber.App) {
	log.Debug().Msg("listening at port: " + configs.Port)
	err := appInstance.Listen(":" + configs.Port)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}
}

func waitForShutdown(appInstance *fiber.App, rdbmsInstance db.RDBMS, inMemoryDbInstance db.InMemoryDB) {
	sigChan := make(chan os.Signal, 1) // Create a channel to listen for OS signals
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan // Block until a signal is received
	log.Debug().Msg("received termination signal, shutting down")
	app.CloseConnections(appInstance, rdbmsInstance, inMemoryDbInstance)
}

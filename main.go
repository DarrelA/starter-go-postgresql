package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/DarrelA/starter-go-postgresql/app"
	pgdb "github.com/DarrelA/starter-go-postgresql/db/pgdb"
	redisDb "github.com/DarrelA/starter-go-postgresql/db/redis"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
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

	go app.StartApp(rdbmsInstance, inMemoryDbInstance)

	// Create a channel to listen for OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan                                               // Block until a signal is received
	app.CloseConnections(rdbmsInstance, inMemoryDbInstance) // Gracefully close the connections

	logFile.Close()
	os.Exit(0)
}

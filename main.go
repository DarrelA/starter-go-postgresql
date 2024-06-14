package main

import (
	"os"

	"github.com/DarrelA/starter-go-postgresql/app"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Create a logger instance with output to a file
	logFile, err := os.Create("/app/logs/app.log")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create log file")
	}
	log.Info().Msg("log file has been created")
	defer logFile.Close()

	// Configure logger to write to both file and console
	log.Logger = zerolog.
		New(zerolog.
			MultiLevelWriter(zerolog.ConsoleWriter{Out: os.Stderr}, logFile)).
		With().
		Caller().
		Timestamp().
		Logger()

	app.StartApp()
}

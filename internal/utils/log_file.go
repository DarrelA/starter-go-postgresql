package utils

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var logFile *os.File

func CreateAppLog() {
	var err error

	// Create a logger instance with output to a file
	logFile, err = os.Create("/app/logs/app.log")
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
}

func GetLogFile() *os.File {
	return logFile
}

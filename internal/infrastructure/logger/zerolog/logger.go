package zerolog

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type ZeroLogger struct{}

func NewZeroLogger(logFile *os.File) *ZeroLogger {
	// Configure logger to write to both file and console
	log.Logger = zerolog.
		New(zerolog.
			MultiLevelWriter(zerolog.ConsoleWriter{Out: os.Stderr}, logFile)).
		With().
		Caller().
		Timestamp().
		Logger()

	return &ZeroLogger{}
}

package zerolog

import (
	"os"

	repo "github.com/DarrelA/starter-go-postgresql/internal/domain/repository"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type ZeroLogger struct{}

func NewZeroLogger(logFile *os.File) repo.DBLogger {
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

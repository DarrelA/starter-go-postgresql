package envs_utils

import (
	"os"
	"os/exec"
	"runtime"

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

func LogCWD() string {
	// Get caller information
	pc, file, line, ok := runtime.Caller(1)
	var callerInfo string
	if ok {
		caller := runtime.FuncForPC(pc)
		callerInfo = caller.Name()
		log.Debug().Msgf("LogCWD() is called by %s (%s:%d)", callerInfo, file, line)
	} else {
		log.Debug().Msg("failed to get caller information")
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Error().Err(err).Msg("error getting current working directory")
	}

	log.Debug().Msgf("@cwd: %s", cwd)
	return cwd
}

func ListFiles() {
	// Execute the `ls` command
	cmd := exec.Command("ls")
	output, err := cmd.Output()
	if err != nil {
		log.Error().Err(err).Msg("error executing ls command")
	}

	// Log the output of `ls` command
	log.Debug().Str("ls_output", string(output)).Msg("List of files")
}

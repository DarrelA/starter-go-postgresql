package logger

import (
	"os"
	"os/exec"
	"runtime"

	"github.com/rs/zerolog/log"
)

const (
	errMsgCreateLogFileError       = "failed to create log file"
	errMsgFailedToGetCallerInfo    = "failed to get caller information"
	errMsgGetWorkingDirectoryError = "error getting current working directory"
	errMsgExecutingCmd             = "error executing ls command"
)

func CreateAppLog(logFilePath string) *os.File {
	// Create a logger instance with output to a file
	logFile, err := os.Create(logFilePath)
	if err != nil {
		log.Error().Err(err).Msg(errMsgCreateLogFileError)
	}
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
		log.Debug().Msg(errMsgFailedToGetCallerInfo)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Error().Err(err).Msg(errMsgGetWorkingDirectoryError)
	}

	log.Debug().Msgf("@cwd: %s", cwd)
	return cwd
}

func ListFiles() {
	// Execute the `ls` command
	cmd := exec.Command("ls")
	output, err := cmd.Output()
	if err != nil {
		log.Error().Err(err).Msg(errMsgExecutingCmd)
	}

	// Log the output of `ls` command
	log.Debug().Str("ls_output", string(output)).Msg("List of files")
}

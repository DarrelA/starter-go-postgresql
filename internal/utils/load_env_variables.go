package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	"github.com/DarrelA/starter-go-postgresql/internal/utils/err_rest"
	"github.com/rs/zerolog/log"
)

func LoadEnvVariableInt(envVar string, target *int) {
	valueStr := os.Getenv(envVar)
	if valueStr == "" {
		errMessage := fmt.Sprintf("%s is not set", envVar)
		err := err_rest.NewInternalServerError(errMessage)
		log.Error().Err(err).Msg("")
		return
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Error().Err(err).Msgf("check JWT config: %s", envVar)
		return
	}
	*target = value
}

func LoadEnvVariableBool(envVar string, target *bool) {
	valueStr := os.Getenv(envVar)
	if valueStr == "" {
		errMessage := fmt.Sprintf("%s is not set", envVar)
		err := err_rest.NewInternalServerError(errMessage)
		log.Error().Err(err).Msg("")
		return
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		log.Error().Err(err).Msgf("check JWT config: %s", envVar)
		return
	}
	*target = value
}

func LoadEnvVariableDuration(envVar string, target *time.Duration) {
	valueStr := os.Getenv(envVar)
	if valueStr == "" {
		errMessage := fmt.Sprintf("%s is not set", envVar)
		err := err_rest.NewInternalServerError(errMessage)
		log.Error().Err(err).Msg("")
		return
	}
	value, err := time.ParseDuration(valueStr)
	if err != nil {
		log.Error().Err(err).Msgf("check JWT config: %s", envVar)
		return
	}
	*target = value
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

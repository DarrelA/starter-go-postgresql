package utils

import (
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

func LoadEnvVariableInt(envVar string, target *int) {
	valueStr := os.Getenv(envVar)
	if valueStr == "" {
		log.Error().Msgf("%s is not set", envVar)
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
		log.Error().Msgf("%s is not set", envVar)
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
		log.Error().Msgf("%s is not set", envVar)
		return
	}
	value, err := time.ParseDuration(valueStr)
	if err != nil {
		log.Error().Err(err).Msgf("check JWT config: %s", envVar)
		return
	}
	*target = value
}

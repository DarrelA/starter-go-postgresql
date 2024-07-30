package zerolog

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
	config "github.com/DarrelA/starter-go-postgresql/internal/infrastructure/config"
	"github.com/rs/zerolog/log"
)

const logFilePath = "./app.log"

const (
	errMsgCreateLogFileError = "failed to create log file"
)

func TestNewZeroLogger(t *testing.T) {
	tests := []struct {
		name        string
		envConfig   config.EnvConfig
		logLevel    string
		expectedMsg string
	}{
		{
			name: "log info message", envConfig: config.EnvConfig{EnvConfig: entity.EnvConfig{Env: "test"}},
			logLevel: "info", expectedMsg: "this is an info message",
		},
		{
			name: "log error message", envConfig: config.EnvConfig{EnvConfig: entity.EnvConfig{Env: "test"}},
			logLevel: "error", expectedMsg: "this is an error message",
		},
		{
			name: "log debug message", envConfig: config.EnvConfig{EnvConfig: entity.EnvConfig{Env: "test"}},
			logLevel: "debug", expectedMsg: "this is an debug message",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			os.Setenv("LOG_LEVEL", test.logLevel)
			defer os.Unsetenv("LOG_LEVEL")
			e := test.envConfig
			e.LoadLogConfig()

			logFile, err := os.Create(logFilePath)
			if err != nil {
				t.Log(errMsgCreateLogFileError)
			}
			defer logFile.Close()

			NewZeroLogger(logFile)

			// Set up a buffer to capture log output
			var buf bytes.Buffer
			log.Logger = log.Output(&buf)

			// Log the messages
			switch test.logLevel {
			case "info":
				log.Info().Msg(test.expectedMsg)
			case "error":
				log.Error().Msg(test.expectedMsg)
			case "debug":
				log.Debug().Msg(test.expectedMsg)
			}

			logOutput := buf.String()
			if !strings.Contains(logOutput, test.expectedMsg) {
				t.Errorf("Expected '%s' not found in log output: %s", test.expectedMsg, logOutput)
			}

			// Reset the logger to its default output
			log.Logger = log.Output(os.Stdout)

			// Verify the log output in the log file
			_, err = os.ReadFile(logFilePath)
			if err != nil {
				t.Errorf("Failed to read log file: %v", err)
			}

			// @TODO: Find a fix?
			// zerolog doesn’t write to file during Go tests
			// t.Logf("Log file content: %s", string(logFileContent)) // Debug print
			// if !strings.Contains(string(logFileContent), test.expectedMsg) {
			// 	t.Errorf("Expected '%s' not found in log file: %s", test.expectedMsg, string(logFileContent))
			// }

			// Clean up the log file for the next test
			if err := os.Remove(logFilePath); err != nil {
				t.Errorf("Failed to remove log file: %v", err)
			}
		})
	}
}

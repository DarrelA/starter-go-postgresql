package zerolog

import (
	"bytes"
	"os"
	"strings"
	"testing"

	config "github.com/DarrelA/starter-go-postgresql/internal/infrastructure/config"
	envLogger "github.com/DarrelA/starter-go-postgresql/internal/infrastructure/logger"
	"github.com/rs/zerolog/log"
)

const logFilePath = "./app.log"

var tests = []struct {
	name        string
	logLevel    string
	expectedMsg string
}{
	{
		name:     "log info message",
		logLevel: "info", expectedMsg: "this is an info message",
	},
	{
		name:     "log error message",
		logLevel: "error", expectedMsg: "this is an error message",
	},
	{
		name:     "log debug message",
		logLevel: "debug", expectedMsg: "this is a debug message",
	},
}

func TestNewZeroLogger(t *testing.T) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := config.ConfigureLogging(test.logLevel); err != nil {
				t.Fatalf("configure logging: %v", err)
			}
			logFile := envLogger.CreateAppLog(logFilePath)
			defer logFile.Close()

			// Initialize the logger to write to the log file
			NewZeroLogger(logFile)
			logMsg(test.logLevel, test.expectedMsg)

			// Verify the log output in the log file
			logFileContent, err := os.ReadFile(logFilePath)
			if err != nil {
				t.Errorf("Failed to read log file: %v", err)
			}

			if !strings.Contains(string(logFileContent), test.expectedMsg) {
				t.Errorf("Expected '%s' not found in log file: %s", test.expectedMsg, string(logFileContent))
			}

			// Clean up the log file for the next test
			if err := os.Remove(logFilePath); err != nil {
				t.Errorf("Failed to remove log file: %v", err)
			}

			// Set up a buffer to capture log output
			var buf bytes.Buffer
			log.Logger = log.Output(&buf)
			logMsg(test.logLevel, test.expectedMsg)
			logOutput := buf.String()
			if !strings.Contains(logOutput, test.expectedMsg) {
				t.Errorf("Expected '%s' not found in log output: %s", test.expectedMsg, logOutput)
			}

			// Reset the logger to its default output
			log.Logger = log.Output(os.Stdout)
		})
	}
}

func logMsg(logLevel string, expectedMsg string) {
	switch logLevel {
	case "info":
		log.Info().Msg(expectedMsg)
	case "error":
		log.Error().Msg(expectedMsg)
	case "debug":
		log.Debug().Msg(expectedMsg)
	}
}

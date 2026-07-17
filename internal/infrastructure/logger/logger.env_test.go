package logger

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestCreateAppLog(t *testing.T) {
	// Create a buffer to capture stdout
	var buf bytes.Buffer

	// Set zerolog to write to the buffer
	log.Logger = log.Output(&buf)
	zerolog.SetGlobalLevel(zerolog.ErrorLevel)

	// Call CreateAppLog
	logFile := CreateAppLog("👾👾👾/invalid/path/to/app.log")

	// Ensure logFile is nil due to the error
	if logFile != nil {
		t.Errorf("Expected nil logFile, got %v", logFile)
	}

	// Capture the log output
	logOutput := buf.String()

	// Check if the string contains the expected error message
	expectedMessage := errMsgCreateLogFileError
	if !strings.Contains(logOutput, expectedMessage) {
		t.Errorf("Expected error message not found in log output: %s", logOutput)
	}

	// Clean up by resetting the logger
	log.Logger = log.Output(os.Stdout)
}

func TestLogCWD(t *testing.T) {
	t.Run("cwd", func(t *testing.T) {
		var buf bytes.Buffer
		log.Logger = log.Output(&buf)
		zerolog.SetGlobalLevel(zerolog.DebugLevel)

		cwd := LogCWD()

		if cwd == "" {
			t.Errorf("Expected non-empty current working directory, got empty string")
		}

		logOutput := buf.String()
		expectedDebugMessage := "@cwd: " + cwd
		if !strings.Contains(logOutput, expectedDebugMessage) {
			t.Errorf("Expected debug message not found in log output: %s", logOutput)
		}

		log.Logger = log.Output(os.Stdout)
	})

	t.Run("CallerInfo", func(t *testing.T) {
		var buf bytes.Buffer
		log.Logger = log.Output(&buf)
		zerolog.SetGlobalLevel(zerolog.DebugLevel)

		LogCWD()

		logOutput := buf.String()
		expectedCallerMessage := "LogCWD() is called by"
		if !strings.Contains(logOutput, expectedCallerMessage) {
			t.Errorf("Expected '%s' in log output: %s", expectedCallerMessage, logOutput)
		}

		log.Logger = log.Output(os.Stdout)
	})
}

func TestListFiles(t *testing.T) {
	var buf bytes.Buffer
	log.Logger = log.Output(&buf)
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	ListFiles()

	logOutput := buf.String()
	expectedDebugMessage := "List of files"
	if !strings.Contains(logOutput, expectedDebugMessage) {
		t.Errorf("Expected '%s' in log output: %s", expectedDebugMessage, logOutput)
	}

	if !strings.Contains(logOutput, "ls_output") {
		t.Errorf("Expected 'ls output' not found in log output: %s", logOutput)
	}

	log.Logger = log.Output(os.Stdout)
}

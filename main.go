package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/DarrelA/starter-go-postgresql/app"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Create a logger instance with output to a file
	logFile, err := os.Create("/app/logs/app.log")
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

	// Use `WaitGroup` when you just need to wait for tasks to complete without exchanging data.
	// Use channels when you need to signal task completion and possibly exchange data.
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		app.CreateDBConnections()
		app.ConfigureAppInstance()
	}()

	wg.Wait()
	go app.StartServer()
	waitForShutdown()

	logFile.Close()
	os.Exit(0)
}

func waitForShutdown() {
	sigChan := make(chan os.Signal, 1) // Create a channel to listen for OS signals
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan // Block until a signal is received
	log.Debug().Msg("received termination signal, shutting down")
	app.CloseConnections()
}

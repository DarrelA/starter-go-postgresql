package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/DarrelA/starter-go-postgresql/app"
	"github.com/DarrelA/starter-go-postgresql/internal/utils"
	"github.com/rs/zerolog/log"
)

func main() {
	utils.CreateAppLog()

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

	utils.GetLogFile().Close()
	os.Exit(0)
}

func waitForShutdown() {
	sigChan := make(chan os.Signal, 1) // Create a channel to listen for OS signals
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan // Block until a signal is received
	log.Debug().Msg("received termination signal, shutting down")
	app.CloseConnections()
}

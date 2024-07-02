package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/DarrelA/starter-go-postgresql/app"
	envs_utils "github.com/DarrelA/starter-go-postgresql/internal/utils/envs"
	"github.com/rs/zerolog/log"
)

func main() {
	envs_utils.CreateAppLog()

	// Use `WaitGroup` when you just need to wait for tasks to complete without exchanging data.
	// Use channels when you need to signal task completion and possibly exchange data.
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		// Functions are running synchronously
		app.CreateDBConnections()
		app.SeedDatabase()
		app.ConfigureAppInstance()
	}()

	wg.Wait()
	go app.StartServer()
	waitForShutdown()

	envs_utils.GetLogFile().Close()
	os.Exit(0)
}

func waitForShutdown() {
	sigChan := make(chan os.Signal, 1) // Create a channel to listen for OS signals
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan // Block until a signal is received
	log.Debug().Msg("received termination signal, shutting down")
	app.CloseConnections()
}

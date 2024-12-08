package utils

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// GracefulShutdown registers cleanup functions and handles graceful shutdown of the application.
func GracefulShutdown(ctx context.Context, cancel context.CancelFunc, cleanupFuncs ...func() error) {
	// Set up a signal handler to catch SIGINT (Ctrl+C) and SIGTERM (termination signals)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Start a goroutine to listen for signals
	go func() {
		// Wait for a termination signal
		<-sigs
		log.Println("Shutting down gracefully...")

		// Call cleanup functions
		for _, cleanup := range cleanupFuncs {
			if err := cleanup(); err != nil {
				log.Printf("Error during cleanup: %v", err)
			}
		}

		// Cancel the context to stop the application
		cancel()
	}()
}

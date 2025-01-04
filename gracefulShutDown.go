package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// handler simulates a request handler with some processing delay.
func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling request: %s %s", r.Method, r.URL.Path)
	time.Sleep(2 * time.Second) // Simulate a long-running request
	fmt.Fprintf(w, "Hello, you've hit %s\n", r.URL.Path)
}

func main() {
	// Create a new HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(handler),
	}

	// Channel to listen for OS signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Channel to signal when the server has shut down gracefully
	shutdownComplete := make(chan struct{})

	// Create a WaitGroup to track ongoing requests
	var wg sync.WaitGroup

	// Wrap the handler to manage ongoing requests using the WaitGroup
	server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wg.Add(1)
		defer wg.Done()
		handler(w, r)
	})

	// Goroutine to listen for shutdown signals
	go func() {
		<-signalChan
		log.Println("Shutdown signal received. Stopping server...")

		// Create a context with a timeout for the shutdown process
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Shutdown the server gracefully
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Error during server shutdown: %v", err)
		}

		// Wait for all ongoing requests to complete
		wg.Wait()

		// Signal that shutdown is complete
		close(shutdownComplete)
	}()

	// Start the server
	log.Println("Starting server on :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}

	// Wait for the shutdown process to complete
	<-shutdownComplete
	log.Println("Server has shut down gracefully.")
}

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	defaultPort           = "8081"
	shutdownTimeout       = 30 * time.Second
	readHeaderTimeout     = 10 * time.Second
	readTimeout           = 15 * time.Second
	writeTimeout          = 15 * time.Second
	idleTimeout           = 60 * time.Second
)

func main() {
	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Create HTTP server
	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		Handler:           setupRouter(),
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting redirect-service on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}

// setupRouter creates and configures the HTTP router
func setupRouter() http.Handler {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", healthHandler)

	// Readiness check endpoint
	mux.HandleFunc("/ready", readyHandler)

	// Placeholder for redirect endpoint (to be implemented)
	mux.HandleFunc("/", placeholderHandler)

	return mux
}

// healthHandler responds to liveness probe requests
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"healthy","service":"redirect-service","timestamp":"%s"}`, time.Now().UTC().Format(time.RFC3339))
}

// readyHandler responds to readiness probe requests
func readyHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Add actual dependency checks (Redis, PostgreSQL)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"ready","service":"redirect-service","timestamp":"%s"}`, time.Now().UTC().Format(time.RFC3339))
}

// placeholderHandler is a temporary handler until the redirect logic is implemented
func placeholderHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, `{"error":"not_implemented","message":"Redirect service is not yet fully implemented"}`)
}

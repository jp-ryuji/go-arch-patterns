package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jp-ryuji/go-arch-patterns/internal/config"
	"github.com/jp-ryuji/go-arch-patterns/internal/di"
	"github.com/jp-ryuji/go-arch-patterns/internal/infrastructure/postgres"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create database client
	client := postgres.NewClient(cfg.DatabaseURL())

	// Create dependency injection container
	container, err := di.NewContainer(client, cfg.GRPCPort, cfg.HTTPPort)
	if err != nil {
		log.Fatalf("Failed to create container: %v", err)
	}

	// Start the server
	log.Println("Starting server...")
	if err := container.HTTPServer.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	// Give the server 5 seconds to shutdown gracefully
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	if err := container.HTTPServer.Stop(); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

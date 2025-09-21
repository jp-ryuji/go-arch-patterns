package main

import (
	"context"
	"log"

	"github.com/jp-ryuji/go-arch-patterns/internal/config"
	"github.com/jp-ryuji/go-arch-patterns/internal/infrastructure/postgres"
)

func main() {
	// Load configuration using Viper
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create Ent client using our NewClient function
	// Note: For migration, we need to connect to localhost when running from host machine
	// But use the docker service name when running from within docker
	databaseUrl := cfg.DatabaseURL()
	client := postgres.NewClient(databaseUrl)
	defer client.Close()

	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	log.Println("Migration completed successfully")
}

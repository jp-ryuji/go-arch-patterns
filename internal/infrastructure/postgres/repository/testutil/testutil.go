//go:build integration

// Package testutil provides shared test utilities for repository integration tests.
package testutil

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/jp-ryuji/go-arch-patterns/internal/infrastructure/postgres"
	"github.com/jp-ryuji/go-arch-patterns/internal/infrastructure/postgres/entgen"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var (
	DBClient *entgen.Client       // shared database client for all repository tests
	Pool     *dockertest.Pool     // shared Docker test pool
	Resource *dockertest.Resource // shared Docker resource
)

// SetupTestEnvironment initializes the shared test environment
func SetupTestEnvironment() error {
	log.Printf("Setting up test environment...")

	// Create a Docker pool
	var err error
	Pool, err = dockertest.NewPool("")
	if err != nil {
		log.Printf("Failed to create Docker pool: %v", err)
		return fmt.Errorf("could not construct pool: %w", err)
	}

	// Ping the Docker daemon to ensure it's running
	log.Printf("Pinging Docker daemon...")
	err = Pool.Client.Ping()
	if err != nil {
		log.Printf("Failed to ping Docker daemon: %v", err)
		return fmt.Errorf("could not connect to Docker: %w", err)
	}

	// Start a PostgreSQL container
	log.Printf("Starting PostgreSQL container...")
	Resource, err = Pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "17",
		Env: []string{
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_USER=user",
			"POSTGRES_DB=dbname",
		},
	}, func(config *docker.HostConfig) {
		// Set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Printf("Failed to start PostgreSQL container: %v", err)
		return fmt.Errorf("could not start resource: %w", err)
	}

	// Set a longer timeout for Docker to start the container
	Pool.MaxWait = 120 * time.Second

	// Wait a bit for the container to start
	log.Printf("Waiting for container to start...")
	time.Sleep(5 * time.Second)

	// Retry connecting to the database until it's ready
	log.Printf("Retrying database connection...")
	if err = Pool.Retry(func() error {
		hostPort := Resource.GetPort("5432/tcp")
		log.Printf("Attempting to connect to database at 127.0.0.1:%s", hostPort)

		// Try with IPv4 localhost explicitly
		databaseUrl := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
			"user",
			"secret",
			fmt.Sprintf("127.0.0.1:%s", hostPort),
			"dbname",
			"disable")

		DBClient = postgres.NewClient(databaseUrl)

		// Just return nil since we don't need to create the uuid extension here
		return nil
	}); err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return fmt.Errorf("could not connect to docker: %w", err)
	}

	// Run database migrations using Ent
	log.Printf("Running database migrations...")
	if err := DBClient.Schema.Create(context.Background()); err != nil {
		log.Printf("Failed to run database migrations: %v", err)
		return fmt.Errorf("could not migrate: %w", err)
	}

	log.Printf("Test environment setup completed successfully")
	return nil
}

// TeardownTestEnvironment cleans up the shared test environment
func TeardownTestEnvironment() error {
	log.Printf("Tearing down test environment...")
	if Pool != nil && Resource != nil {
		log.Printf("Purging Docker resource...")
		if err := Pool.Purge(Resource); err != nil {
			log.Printf("Warning: could not purge resource: %v", err)
			// Don't return the error, just log it
		}
		log.Printf("Successfully purged Docker resource")
	} else {
		log.Printf("Pool or Resource is nil, skipping purge")
	}
	log.Printf("Finished tearing down test environment")
	return nil
}

// SkipIfShort skips the test if running in short mode
func SkipIfShort(t *testing.T) {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping integration test")
	}
}

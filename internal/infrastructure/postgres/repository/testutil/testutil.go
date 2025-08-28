//go:build integration

// Package testutil provides shared test utilities for repository integration tests.
package testutil

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/jp-ryuji/go-sample/internal/infrastructure/postgres"
	"github.com/jp-ryuji/go-sample/internal/infrastructure/postgres/dbmodel"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var (
	DBClient *postgres.Client     // shared database client for all repository tests
	Pool     *dockertest.Pool     // shared Docker test pool
	Resource *dockertest.Resource // shared Docker resource
)

// allModels is a slice of all GORM models used in the application
// When adding new models, they should be added to this slice
var allModels = []interface{}{
	&dbmodel.Car{},
	&dbmodel.Company{},
	&dbmodel.Tenant{},
}

// SetupTestEnvironment initializes the shared test environment
func SetupTestEnvironment() error {
	// Create a Docker pool
	var err error
	Pool, err = dockertest.NewPool("")
	if err != nil {
		return fmt.Errorf("could not construct pool: %w", err)
	}

	// Ping the Docker daemon to ensure it's running
	err = Pool.Client.Ping()
	if err != nil {
		return fmt.Errorf("could not connect to Docker: %w", err)
	}

	// Start a PostgreSQL container
	Resource, err = Pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "13",
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
		return fmt.Errorf("could not start resource: %w", err)
	}

	// Set a longer timeout for Docker to start the container
	Pool.MaxWait = 120 * time.Second

	// Wait a bit for the container to start
	time.Sleep(5 * time.Second)

	// Retry connecting to the database until it's ready
	if err = Pool.Retry(func() error {
		hostPort := Resource.GetPort("5432/tcp")
		log.Printf("Attempting to connect to database at 127.0.0.1:%s", hostPort)

		// Try with IPv4 localhost explicitly
		DBClient = postgres.NewClient(
			fmt.Sprintf("127.0.0.1:%s", hostPort),
			"user",
			"secret",
			"dbname",
			false,
		)
		return DBClient.DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error
	}); err != nil {
		return fmt.Errorf("could not connect to docker: %w", err)
	}

	// Run database migrations using the actual GORM models
	if err := DBClient.DB.AutoMigrate(allModels...); err != nil {
		return fmt.Errorf("could not migrate: %w", err)
	}

	return nil
}

// TeardownTestEnvironment cleans up the shared test environment
func TeardownTestEnvironment() error {
	if err := Pool.Purge(Resource); err != nil {
		return fmt.Errorf("could not purge resource: %w", err)
	}
	return nil
}

// SkipIfShort skips the test if running in short mode
func SkipIfShort(t *testing.T) {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping integration test")
	}
}

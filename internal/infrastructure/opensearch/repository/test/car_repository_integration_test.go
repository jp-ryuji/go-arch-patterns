package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jp-ryuji/go-sample/internal/domain/model"
	"github.com/jp-ryuji/go-sample/internal/infrastructure/opensearch/repository"
	opensearch "github.com/opensearch-project/opensearch-go/v4"
	opensearchapi "github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCarRepositoryIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Create a Docker pool
	pool, err := dockertest.NewPool("")
	require.NoError(t, err)

	// Run OpenSearch container
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "opensearchproject/opensearch",
		Tag:        "2.11.0", // Use a specific version for consistency
		Env: []string{
			"discovery.type=single-node",
			"OPENSEARCH_JAVA_OPTS=-Xms512m -Xmx512m",
			"DISABLE_SECURITY_PLUGIN=true", // Disable security for testing
			"DISABLE_INSTALL_DEMO_CONFIG=true",
		},
		ExposedPorts: []string{"9200/tcp"},
	}, func(config *docker.HostConfig) {
		// Set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	require.NoError(t, err)

	// Clean up the container after the test
	t.Cleanup(func() {
		require.NoError(t, pool.Purge(resource))
	})

	// Get the host port
	hostPort := resource.GetPort("9200/tcp")

	// Create OpenSearch client configuration
	cfg := opensearch.Config{
		Addresses: []string{fmt.Sprintf("http://localhost:%s", hostPort)},
	}

	// Test connection
	err = pool.Retry(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		apiClient, apiErr := opensearchapi.NewClient(opensearchapi.Config{Client: cfg})
		if apiErr != nil {
			return apiErr
		}
		_, pingErr := apiClient.Ping(ctx, nil)
		return pingErr
	})
	require.NoError(t, err)

	// Create repository
	repo := repository.NewCarRepository(cfg, "test-cars")

	t.Run("Create and Get Car", func(t *testing.T) {
		ctx := context.Background()
		now := time.Now()

		// Create a car
		car := model.NewCar("tenant-123", "Toyota Prius", now)
		car.ID = "test-car-1"

		// Create the car in OpenSearch
		err := repo.Create(ctx, car)
		assert.NoError(t, err)

		// Get the car from OpenSearch
		retrievedCar, err := repo.Get(ctx, car.ID)
		assert.NoError(t, err)
		assert.NotNil(t, retrievedCar)
		assert.Equal(t, car.ID, retrievedCar.ID)
		assert.Equal(t, car.TenantID, retrievedCar.TenantID)
		assert.Equal(t, car.Model, retrievedCar.Model)
	})

	t.Run("Update Car", func(t *testing.T) {
		ctx := context.Background()
		now := time.Now()

		// Create a car
		car := model.NewCar("tenant-123", "Toyota Prius", now)
		car.ID = "test-car-2"

		// Create the car in OpenSearch
		err := repo.Create(ctx, car)
		assert.NoError(t, err)

		// Update the car
		car.Model = "Honda Civic"
		car.UpdatedAt = time.Now()

		// Update the car in OpenSearch
		err = repo.Update(ctx, car)
		assert.NoError(t, err)

		// Get the car from OpenSearch to verify the update
		retrievedCar, err := repo.Get(ctx, car.ID)
		assert.NoError(t, err)
		assert.NotNil(t, retrievedCar)
		assert.Equal(t, car.ID, retrievedCar.ID)
		assert.Equal(t, car.TenantID, retrievedCar.TenantID)
		assert.Equal(t, "Honda Civic", retrievedCar.Model)
	})

	t.Run("Delete Car", func(t *testing.T) {
		ctx := context.Background()
		now := time.Now()

		// Create a car
		car := model.NewCar("tenant-123", "Toyota Prius", now)
		car.ID = "test-car-3"

		// Create the car in OpenSearch
		err := repo.Create(ctx, car)
		assert.NoError(t, err)

		// Delete the car from OpenSearch
		err = repo.Delete(ctx, car.ID)
		assert.NoError(t, err)

		// Try to get the deleted car - should return an error
		_, err = repo.Get(ctx, car.ID)
		assert.Error(t, err)
		// The error might not contain "car not found" exactly, but should indicate the car wasn't found
		// We'll check that it's an error and not nil
	})
}

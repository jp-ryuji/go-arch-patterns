package client

import (
	"context"
	"fmt"
	"os"
	"time"

	opensearch "github.com/opensearch-project/opensearch-go/v4"
	opensearchapi "github.com/opensearch-project/opensearch-go/v4/opensearchapi"
)

// NewClient creates a new OpenSearch client
func NewClient() (*opensearch.Client, error) {
	cfg := opensearch.Config{
		Addresses: []string{
			fmt.Sprintf("http://localhost:%s", getEnv("OPENSEARCH_PORT_EXTERNAL", "9200")),
		},
		Username: "admin",
		Password: getEnv("OPENSEARCH_INITIAL_ADMIN_PASSWORD", "EX2R3L(,M-tR"),
	}

	// Retry connection a few times since OpenSearch might not be ready immediately
	var client *opensearch.Client
	var err error

	for i := 0; i < 10; i++ {
		client, err = opensearch.NewClient(cfg)
		if err == nil {
			// Test the connection
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			apiClient, apiErr := opensearchapi.NewClient(opensearchapi.Config{Client: cfg})
			if apiErr == nil {
				_, pingErr := apiClient.Ping(ctx, nil)
				cancel()
				if pingErr == nil {
					break
				}
			} else {
				cancel()
			}
		}
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create OpenSearch client: %w", err)
	}

	return client, nil
}

// getEnv returns the value of an environment variable or a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetDefaultIndex returns the default index name for cars
func GetDefaultIndex() string {
	return getEnv("OPENSEARCH_CAR_INDEX", "cars")
}

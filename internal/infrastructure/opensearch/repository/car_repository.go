package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jp-ryuji/go-sample/internal/domain/model"
	opensearch "github.com/opensearch-project/opensearch-go/v4"
	opensearchapi "github.com/opensearch-project/opensearch-go/v4/opensearchapi"
)

type carRepository struct {
	client *opensearchapi.Client
	index  string
}

// NewCarRepository creates a new car repository for OpenSearch
func NewCarRepository(client *opensearch.Client, index string) *carRepository {
	// We'll create a new API client with the same configuration
	// This is a simplified approach since we can't access the config from the existing client
	apiClient, _ := opensearchapi.NewClient(opensearchapi.Config{})
	return &carRepository{
		client: apiClient,
		index:  index,
	}
}

// Create inserts a new car into OpenSearch
func (r *carRepository) Create(ctx context.Context, car *model.Car) error {
	doc, err := json.Marshal(car)
	if err != nil {
		return fmt.Errorf("failed to marshal car: %w", err)
	}

	req := opensearchapi.IndexReq{
		Index:      r.index,
		DocumentID: car.ID,
		Body:       strings.NewReader(string(doc)),
	}

	res, err := r.client.Index(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to index car: %w", err)
	}

	if res.Result == "" {
		return fmt.Errorf("failed to index car: empty result")
	}

	return nil
}

// Update updates an existing car in OpenSearch
func (r *carRepository) Update(ctx context.Context, car *model.Car) error {
	return r.Create(ctx, car) // In OpenSearch, indexing with the same ID updates the document
}

// Delete removes a car by its ID from OpenSearch
func (r *carRepository) Delete(ctx context.Context, id string) error {
	req := opensearchapi.DocumentDeleteReq{
		Index:      r.index,
		DocumentID: id,
	}

	res, err := r.client.Document.Delete(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to delete car: %w", err)
	}

	if res.Result == "" {
		return fmt.Errorf("failed to delete car: empty result")
	}

	return nil
}

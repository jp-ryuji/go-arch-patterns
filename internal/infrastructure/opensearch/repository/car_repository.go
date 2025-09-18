package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jp-ryuji/go-arch-patterns/internal/domain/model"
	opensearch "github.com/opensearch-project/opensearch-go/v4"
	opensearchapi "github.com/opensearch-project/opensearch-go/v4/opensearchapi"
)

type CarRepository struct {
	config opensearch.Config
	index  string
}

// NewCarRepository creates a new car repository for OpenSearch
func NewCarRepository(config opensearch.Config, index string) *CarRepository {
	return &CarRepository{
		config: config,
		index:  index,
	}
}

// getClient creates a new API client with the stored configuration
func (r *CarRepository) getClient() (*opensearchapi.Client, error) {
	apiClient, err := opensearchapi.NewClient(opensearchapi.Config{Client: r.config})
	if err != nil {
		return nil, fmt.Errorf("failed to create API client: %w", err)
	}
	return apiClient, nil
}

// Create inserts a new car into OpenSearch
func (r *CarRepository) Create(ctx context.Context, car *model.Car) error {
	doc, err := json.Marshal(car)
	if err != nil {
		return fmt.Errorf("failed to marshal car: %w", err)
	}

	client, err := r.getClient()
	if err != nil {
		return err
	}

	req := opensearchapi.IndexReq{
		Index:      r.index,
		DocumentID: car.ID,
		Body:       strings.NewReader(string(doc)),
	}

	res, err := client.Index(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to index car: %w", err)
	}

	if res.Result == "" {
		return fmt.Errorf("failed to index car: empty result")
	}

	return nil
}

// Get retrieves a car by its ID from OpenSearch
func (r *CarRepository) Get(ctx context.Context, id string) (*model.Car, error) {
	client, err := r.getClient()
	if err != nil {
		return nil, err
	}

	req := opensearchapi.DocumentGetReq{
		Index:      r.index,
		DocumentID: id,
	}

	res, err := client.Document.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get car: %w", err)
	}

	if !res.Found {
		return nil, fmt.Errorf("car not found")
	}

	// Parse the car from the response
	var car model.Car
	if err := json.Unmarshal(res.Source, &car); err != nil {
		return nil, fmt.Errorf("failed to unmarshal car: %w", err)
	}

	return &car, nil
}

// Update updates an existing car in OpenSearch
func (r *CarRepository) Update(ctx context.Context, car *model.Car) error {
	return r.Create(ctx, car) // In OpenSearch, indexing with the same ID updates the document
}

// Delete removes a car by its ID from OpenSearch
func (r *CarRepository) Delete(ctx context.Context, id string) error {
	client, err := r.getClient()
	if err != nil {
		return err
	}

	req := opensearchapi.DocumentDeleteReq{
		Index:      r.index,
		DocumentID: id,
	}

	res, err := client.Document.Delete(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to delete car: %w", err)
	}

	if res.Result == "" {
		return fmt.Errorf("failed to delete car: empty result")
	}

	return nil
}

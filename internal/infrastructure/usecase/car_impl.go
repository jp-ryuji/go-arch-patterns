package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/jp-ryuji/go-sample/internal/domain/model"
	"github.com/jp-ryuji/go-sample/internal/domain/repository"
	"github.com/jp-ryuji/go-sample/internal/infrastructure/usecase/input"
)

// carUsecase implements CarUsecase interface
type carUsecase struct {
	carRepo        repository.CarRepository
	opensearchRepo repository.OpenSearchCarRepository
}

// NewCarUsecase creates a new car usecase
func NewCarUsecase(carRepo repository.CarRepository, opensearchRepo repository.OpenSearchCarRepository) CarUsecase {
	return &carUsecase{
		carRepo:        carRepo,
		opensearchRepo: opensearchRepo,
	}
}

// Register registers a new car with saga pattern
func (uc *carUsecase) Register(ctx context.Context, input input.RegisterCarInput) (*model.Car, error) {
	// Validate input
	if err := Validate(input); err != nil {
		return nil, err
	}

	now := time.Now()
	car := model.NewCar(input.TenantID, input.Model, now)

	// Step 1: Save to PostgreSQL
	if err := uc.carRepo.Create(ctx, car); err != nil {
		return nil, fmt.Errorf("failed to create car in database: %w", err)
	}

	// Step 2: Save to OpenSearch (saga pattern)
	if err := uc.opensearchRepo.Create(ctx, car); err != nil {
		// Compensate: Rollback PostgreSQL insert
		_ = uc.carRepo.Delete(ctx, car.ID) // Best effort rollback
		return nil, fmt.Errorf("failed to create car in opensearch: %w", err)
	}

	return car, nil
}

// Update updates an existing car with saga pattern
func (uc *carUsecase) Update(ctx context.Context, input input.UpdateCarInput) (*model.Car, error) {
	// Validate input
	if err := Validate(input); err != nil {
		return nil, err
	}

	// Get the existing car from database
	existingCar, err := uc.carRepo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing car: %w", err)
	}

	// Update the car properties
	existingCar.TenantID = input.TenantID
	existingCar.Model = input.Model
	existingCar.UpdatedAt = time.Now()

	// Step 1: Update in PostgreSQL
	if err := uc.carRepo.Update(ctx, existingCar); err != nil {
		return nil, fmt.Errorf("failed to update car in database: %w", err)
	}

	// Step 2: Update in OpenSearch (saga pattern)
	if err := uc.opensearchRepo.Update(ctx, existingCar); err != nil {
		// Compensate: Rollback PostgreSQL update by reverting to original values
		// In a real implementation, you would save the original state before updating
		// For this example, we'll just return the error
		return nil, fmt.Errorf("failed to update car in opensearch: %w", err)
	}

	return existingCar, nil
}

// GetByID retrieves a car by its ID
func (uc *carUsecase) GetByID(ctx context.Context, input input.GetCarByIDInput) (*model.Car, error) {
	// Validate input
	if err := Validate(input); err != nil {
		return nil, err
	}

	return uc.carRepo.GetByID(ctx, input.ID)
}

// GetByIDWithTenant retrieves a car by its ID along with its tenant information
func (uc *carUsecase) GetByIDWithTenant(ctx context.Context, input input.GetCarByIDInput) (*model.Car, error) {
	// Validate input
	if err := Validate(input); err != nil {
		return nil, err
	}

	return uc.carRepo.GetByIDWithTenant(ctx, input.ID)
}

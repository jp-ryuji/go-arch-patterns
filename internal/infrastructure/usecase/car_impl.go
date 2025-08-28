package usecase

import (
	"context"
	"time"

	"github.com/jp-ryuji/go-sample/internal/domain/model"
	"github.com/jp-ryuji/go-sample/internal/domain/repository"
	"github.com/jp-ryuji/go-sample/internal/infrastructure/usecase/input"
)

// carUsecase implements CarUsecase interface
type carUsecase struct {
	carRepo repository.CarRepository
}

// NewCarUsecase creates a new car usecase
func NewCarUsecase(carRepo repository.CarRepository) CarUsecase {
	return &carUsecase{
		carRepo: carRepo,
	}
}

// Register registers a new car
func (uc *carUsecase) Register(ctx context.Context, input input.RegisterCarInput) (*model.Car, error) {
	// Validate input
	if err := Validate(input); err != nil {
		return nil, err
	}

	now := time.Now()
	car := model.NewCar(input.TenantID, input.Model, now)

	if err := uc.carRepo.Create(ctx, car); err != nil {
		return nil, err
	}

	return car, nil
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

package repository

import (
	"context"

	"github.com/jp-ryuji/go-sample/internal/domain/model"
	"github.com/jp-ryuji/go-sample/internal/domain/repository"
	"github.com/jp-ryuji/go-sample/internal/infrastructure/postgres/dbmodel"
	"github.com/jp-ryuji/go-sample/internal/infrastructure/postgres/query"
)

type carRepository struct{}

// NewCarRepository creates a new car repository
func NewCarRepository() repository.CarRepository {
	return &carRepository{}
}

// Create inserts a new car into the database
func (r *carRepository) Create(ctx context.Context, car *model.Car) error {
	carDB := dbmodel.FromDomain(car)
	return query.Car.WithContext(ctx).Create(carDB)
}

// GetByID retrieves a car by its ID
func (r *carRepository) GetByID(ctx context.Context, id string) (*model.Car, error) {
	carDB, err := query.Car.WithContext(ctx).Where(query.Car.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return carDB.ToDomain(), nil
}

// Update updates an existing car
func (r *carRepository) Update(ctx context.Context, car *model.Car) error {
	carDB := dbmodel.FromDomain(car)
	return query.Car.WithContext(ctx).Save(carDB)
}

// Delete removes a car by its ID
func (r *carRepository) Delete(ctx context.Context, id string) error {
	_, err := query.Car.WithContext(ctx).Where(query.Car.ID.Eq(id)).Delete()
	return err
}

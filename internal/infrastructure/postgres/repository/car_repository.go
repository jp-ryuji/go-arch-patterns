package repository

import (
	"context"

	"github.com/jp-ryuji/go-sample/internal/domain/model"
	"github.com/jp-ryuji/go-sample/internal/domain/repository"
	"github.com/jp-ryuji/go-sample/internal/infrastructure/postgres/dbmodel"
	"gorm.io/gorm"
)

type carRepository struct {
	db *gorm.DB
}

// NewCarRepository creates a new car repository
func NewCarRepository(db *gorm.DB) repository.CarRepository {
	return &carRepository{db: db}
}

// Create inserts a new car into the database
func (r *carRepository) Create(ctx context.Context, car *model.Car) error {
	carDB := dbmodel.FromDomain(car)
	return r.db.WithContext(ctx).Create(carDB).Error
}

// GetByID retrieves a car by its ID
func (r *carRepository) GetByID(ctx context.Context, id string) (*model.Car, error) {
	var carDB dbmodel.Car
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&carDB).Error; err != nil {
		return nil, err
	}
	return carDB.ToDomain(), nil
}

// Update updates an existing car
func (r *carRepository) Update(ctx context.Context, car *model.Car) error {
	carDB := dbmodel.FromDomain(car)
	return r.db.WithContext(ctx).Save(carDB).Error
}

// Delete removes a car by its ID
func (r *carRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&dbmodel.Car{}, "id = ?", id).Error
}

package repository

import (
	"context"

	"github.com/jp-ryuji/go-sample/internal/domain/model"
)

type CarRepository interface {
	Create(ctx context.Context, car *model.Car) error
	GetByID(ctx context.Context, id string) (*model.Car, error)
	Update(ctx context.Context, car *model.Car) error
	Delete(ctx context.Context, id string) error
}

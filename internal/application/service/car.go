package service

import (
	"context"

	"github.com/jp-ryuji/go-arch-patterns/internal/application/input"
	"github.com/jp-ryuji/go-arch-patterns/internal/application/output"
)

// CarService defines the interface for car-related business logic
//
//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_service
type CarService interface {
	Create(ctx context.Context, input input.CreateCar) (*output.CreateCar, error)
	GetByID(ctx context.Context, input input.GetCarByID) (*output.GetCar, error)
	GetByIDWithTenant(ctx context.Context, input input.GetCarByID) (*output.GetCar, error)
	List(ctx context.Context, input input.ListCars) (*output.ListCars, error)
}

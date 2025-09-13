package usecase

import (
	"context"

	"github.com/jp-ryuji/go-ddd/internal/domain/model"
	"github.com/jp-ryuji/go-ddd/internal/infrastructure/usecase/input"
)

// CarUsecase defines the interface for car-related business logic
//
//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_usecase
type CarUsecase interface {
	Register(ctx context.Context, input input.RegisterCarInput) (*model.Car, error)
	GetByID(ctx context.Context, input input.GetCarByIDInput) (*model.Car, error)
	GetByIDWithTenant(ctx context.Context, input input.GetCarByIDInput) (*model.Car, error)
}

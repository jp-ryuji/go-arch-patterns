package service

import (
	"context"

	"github.com/jp-ryuji/go-arch-patterns/internal/application/dto"
	"github.com/jp-ryuji/go-arch-patterns/internal/domain/entity"
)

// CarService defines the interface for car-related business logic
//
//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_service
type CarService interface {
	Register(ctx context.Context, input dto.RegisterCarInput) (*entity.Car, error)
	GetByID(ctx context.Context, input dto.GetCarByIDInput) (*entity.Car, error)
	GetByIDWithTenant(ctx context.Context, input dto.GetCarByIDInput) (*entity.Car, error)
}

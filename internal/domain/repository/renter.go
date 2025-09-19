package repository

import (
	"context"

	"github.com/jp-ryuji/go-arch-patterns/internal/domain/entity"
)

//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_repository
type RenterRepository interface {
	Create(ctx context.Context, renter *entity.Renter) error
	GetByID(ctx context.Context, id string) (*entity.Renter, error)
	Update(ctx context.Context, renter *entity.Renter) error
	Delete(ctx context.Context, id string) error
}

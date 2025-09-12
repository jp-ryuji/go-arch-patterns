package repository

import (
	"context"

	"github.com/jp-ryuji/go-sample/internal/domain/model"
)

//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_repository
type RenterRepository interface {
	Create(ctx context.Context, renter *model.Renter) error
	GetByID(ctx context.Context, id string) (*model.Renter, error)
	Update(ctx context.Context, renter *model.Renter) error
	Delete(ctx context.Context, id string) error
}

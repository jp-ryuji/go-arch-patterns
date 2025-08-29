package repository

import (
	"context"

	"github.com/jp-ryuji/go-sample/internal/domain/model"
)

//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_repository
type OpenSearchCarRepository interface {
	Create(ctx context.Context, car *model.Car) error
	Get(ctx context.Context, id string) (*model.Car, error)
	Update(ctx context.Context, car *model.Car) error
	Delete(ctx context.Context, id string) error
}

package repository

import (
	"context"

	"github.com/jp-ryuji/go-arch-patterns/internal/domain/model"
	"github.com/jp-ryuji/go-arch-patterns/internal/infrastructure/postgres/entgen"
)

//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_repository
type CarRepository interface {
	Create(ctx context.Context, car *model.Car) error
	CreateInTx(ctx context.Context, tx *entgen.Tx, car *model.Car) error
	GetByID(ctx context.Context, id string) (*model.Car, error)
	GetByIDWithTenant(ctx context.Context, id string) (*model.Car, error)
	Update(ctx context.Context, car *model.Car) error
	UpdateInTx(ctx context.Context, tx *entgen.Tx, car *model.Car) error
	Delete(ctx context.Context, id string) error
	DeleteInTx(ctx context.Context, tx *entgen.Tx, id string) error
}

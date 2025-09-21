package repository

import (
	"context"

	"github.com/jp-ryuji/go-arch-patterns/internal/domain/entity"
	"github.com/jp-ryuji/go-arch-patterns/internal/infrastructure/postgres/entgen"
)

// CarLoadOptions defines options for loading related entities
type CarLoadOptions struct {
	WithTenant  bool
	WithRentals bool
}

//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_repository
type CarRepository interface {
	Create(ctx context.Context, car *entity.Car) error
	CreateInTx(ctx context.Context, tx *entgen.Tx, car *entity.Car) error
	GetByID(ctx context.Context, id string) (*entity.Car, error)
	GetByIDWithTenant(ctx context.Context, id string) (*entity.Car, error)
	ListByTenant(ctx context.Context, tenantID string, limit int, offset int) (*entity.Cars, error)
	ListByTenantWithOptions(ctx context.Context, tenantID string, limit int, offset int, opts ...CarLoadOptions) (*entity.Cars, error)
	Update(ctx context.Context, car *entity.Car) error
	UpdateInTx(ctx context.Context, tx *entgen.Tx, car *entity.Car) error
	Delete(ctx context.Context, id string) error
	DeleteInTx(ctx context.Context, tx *entgen.Tx, id string) error
}

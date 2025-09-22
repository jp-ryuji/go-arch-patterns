package repository

import (
	"context"

	"github.com/jp-ryuji/go-arch-patterns/internal/domain/entity"
)

//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_repository
type TenantRepository interface {
	Create(ctx context.Context, tenant *entity.Tenant) error
	GetByID(ctx context.Context, id string) (*entity.Tenant, error)
	GetByCode(ctx context.Context, code string) (*entity.Tenant, error)
	GetByIDWithCars(ctx context.Context, id string) (*entity.Tenant, error)
	Update(ctx context.Context, tenant *entity.Tenant) error
	Delete(ctx context.Context, id string) error
}

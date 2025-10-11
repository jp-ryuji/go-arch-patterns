package rental

import (
	"context"

	sharedrepository "github.com/jp-ryuji/go-arch-patterns/internal/domain/repository"
)

//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=mocks/repository.go -package=mockrental

type CarLoadOptions struct {
	WithTenant  bool
	WithRentals bool
}

type CarRepository interface {
	Create(ctx context.Context, car *Car) error
	CreateInTx(ctx context.Context, tx sharedrepository.Transaction, car *Car) error
	GetByID(ctx context.Context, id string) (*Car, error)
	GetByIDWithTenant(ctx context.Context, id string) (*Car, error)
	ListByTenant(ctx context.Context, tenantID string, limit int, offset int) ([]*Car, string, int32, error)
	ListByTenantWithOptions(ctx context.Context, tenantID string, limit int, offset int, opts ...CarLoadOptions) ([]*Car, string, int32, error)
	Update(ctx context.Context, car *Car) error
	UpdateInTx(ctx context.Context, tx sharedrepository.Transaction, car *Car) error
	Delete(ctx context.Context, id string) error
	DeleteInTx(ctx context.Context, tx sharedrepository.Transaction, id string) error
}

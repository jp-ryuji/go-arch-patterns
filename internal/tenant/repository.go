package tenant

import "context"

//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=mocks/repository.go -package=mocktenant

type Repository interface {
	Create(ctx context.Context, tenant *Tenant) error
	GetByID(ctx context.Context, id string) (*Tenant, error)
	GetByCode(ctx context.Context, code string) (*Tenant, error)
	Update(ctx context.Context, tenant *Tenant) error
	Delete(ctx context.Context, id string) error
}

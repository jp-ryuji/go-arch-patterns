package renter

import "context"

//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=mocks/repository.go -package=mockrenter

type RenterRepository interface {
	Create(ctx context.Context, renter *Renter) error
	GetByID(ctx context.Context, id string) (*Renter, error)
	Update(ctx context.Context, renter *Renter) error
	Delete(ctx context.Context, id string) error
}

type CompanyRepository interface {
	Create(ctx context.Context, company *Company) error
	GetByID(ctx context.Context, id string) (*Company, error)
	Update(ctx context.Context, company *Company) error
	Delete(ctx context.Context, id string) error
}

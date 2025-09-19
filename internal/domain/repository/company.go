package repository

import (
	"context"

	"github.com/jp-ryuji/go-arch-patterns/internal/domain/entity"
)

//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_repository
type CompanyRepository interface {
	Create(ctx context.Context, company *entity.Company) error
	GetByID(ctx context.Context, id string) (*entity.Company, error)
	Update(ctx context.Context, company *entity.Company) error
	Delete(ctx context.Context, id string) error
}

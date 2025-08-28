package repository

import (
	"context"

	"github.com/jp-ryuji/go-sample/internal/domain/model"
)

//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_repository
type CompanyRepository interface {
	Create(ctx context.Context, company *model.Company) error
	GetByID(ctx context.Context, id string) (*model.Company, error)
	Update(ctx context.Context, company *model.Company) error
	Delete(ctx context.Context, id string) error
}

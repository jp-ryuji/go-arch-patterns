package repository

import (
	"context"

	"github.com/jp-ryuji/go-sample/internal/domain/model"
	"github.com/jp-ryuji/go-sample/internal/domain/repository"
	"github.com/jp-ryuji/go-sample/internal/infrastructure/postgres/dbmodel"
	"github.com/jp-ryuji/go-sample/internal/infrastructure/postgres/query"
)

type companyRepository struct{}

// NewCompanyRepository creates a new company repository
func NewCompanyRepository() repository.CompanyRepository {
	return &companyRepository{}
}

// Create inserts a new company into the database
func (r *companyRepository) Create(ctx context.Context, company *model.Company) error {
	companyDB := dbmodel.FromDomainCompany(company)
	return query.Company.WithContext(ctx).Create(companyDB)
}

// GetByID retrieves a company by its ID
func (r *companyRepository) GetByID(ctx context.Context, id string) (*model.Company, error) {
	companyDB, err := query.Company.WithContext(ctx).Where(query.Company.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return companyDB.ToDomain(), nil
}

// Update updates an existing company
func (r *companyRepository) Update(ctx context.Context, company *model.Company) error {
	companyDB := dbmodel.FromDomainCompany(company)
	return query.Company.WithContext(ctx).Save(companyDB)
}

// Delete removes a company by its ID
func (r *companyRepository) Delete(ctx context.Context, id string) error {
	_, err := query.Company.WithContext(ctx).Where(query.Company.ID.Eq(id)).Delete()
	return err
}

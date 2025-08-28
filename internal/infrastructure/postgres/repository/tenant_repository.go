package repository

import (
	"context"

	"github.com/jp-ryuji/go-sample/internal/domain/model"
	"github.com/jp-ryuji/go-sample/internal/domain/repository"
	"github.com/jp-ryuji/go-sample/internal/infrastructure/postgres/dbmodel"
	"github.com/jp-ryuji/go-sample/internal/infrastructure/postgres/query"
)

type tenantRepository struct{}

// NewTenantRepository creates a new tenant repository
func NewTenantRepository() repository.TenantRepository {
	return &tenantRepository{}
}

// Create inserts a new tenant into the database
func (r *tenantRepository) Create(ctx context.Context, tenant *model.Tenant) error {
	tenantDB := dbmodel.FromDomainTenant(tenant)
	return query.Tenant.WithContext(ctx).Create(tenantDB)
}

// GetByID retrieves a tenant by its ID
func (r *tenantRepository) GetByID(ctx context.Context, id string) (*model.Tenant, error) {
	tenantDB, err := query.Tenant.WithContext(ctx).Where(query.Tenant.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return tenantDB.ToDomain(), nil
}

// GetByIDWithCars retrieves a tenant by its ID along with its associated cars
func (r *tenantRepository) GetByIDWithCars(ctx context.Context, id string) (*model.Tenant, error) {
	tenantDB, err := query.Tenant.WithContext(ctx).Preload(query.Tenant.Cars).Where(query.Tenant.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return tenantDB.ToDomain(dbmodel.TenantLoadOptions{WithCars: true}), nil
}

// Update updates an existing tenant
func (r *tenantRepository) Update(ctx context.Context, tenant *model.Tenant) error {
	tenantDB := dbmodel.FromDomainTenant(tenant)
	return query.Tenant.WithContext(ctx).Save(tenantDB)
}

// Delete removes a tenant by its ID
func (r *tenantRepository) Delete(ctx context.Context, id string) error {
	_, err := query.Tenant.WithContext(ctx).Where(query.Tenant.ID.Eq(id)).Delete()
	return err
}

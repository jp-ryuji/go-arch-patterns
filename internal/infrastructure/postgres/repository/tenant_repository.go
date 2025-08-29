package repository

import (
	"context"

	"github.com/jp-ryuji/go-sample/internal/domain/model"
	"github.com/jp-ryuji/go-sample/internal/domain/repository"
	"github.com/jp-ryuji/go-sample/internal/infrastructure/postgres/dbmodel"
	"github.com/jp-ryuji/go-sample/internal/infrastructure/postgres/entgen"
	tenant "github.com/jp-ryuji/go-sample/internal/infrastructure/postgres/entgen/tenant"
)

type tenantRepository struct {
	client *entgen.Client
}

// NewTenantRepository creates a new tenant repository
func NewTenantRepository(client *entgen.Client) repository.TenantRepository {
	return &tenantRepository{client: client}
}

// Create inserts a new tenant into the database
func (r *tenantRepository) Create(ctx context.Context, tenant *model.Tenant) error {
	_, err := r.client.Tenant.
		Create().
		SetID(tenant.ID).
		SetCode(tenant.Code).
		Save(ctx)
	return err
}

// GetByID retrieves a tenant by its ID
func (r *tenantRepository) GetByID(ctx context.Context, id string) (*model.Tenant, error) {
	tenantDB, err := r.client.Tenant.
		Query().
		Where(tenant.ID(id)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	// Convert Ent model to dbmodel and then to domain model
	dbModel := &dbmodel.Tenant{
		ID:        tenantDB.ID,
		Code:      tenantDB.Code,
		CreatedAt: tenantDB.CreatedAt,
		UpdatedAt: tenantDB.UpdatedAt,
	}
	return dbModel.ToDomain(), nil
}

// GetByIDWithCars retrieves a tenant by its ID along with its associated cars
func (r *tenantRepository) GetByIDWithCars(ctx context.Context, id string) (*model.Tenant, error) {
	tenantDB, err := r.client.Tenant.
		Query().
		Where(tenant.ID(id)).
		WithCars().
		Only(ctx)
	if err != nil {
		return nil, err
	}

	// Convert Ent model to dbmodel
	dbModel := &dbmodel.Tenant{
		ID:        tenantDB.ID,
		Code:      tenantDB.Code,
		CreatedAt: tenantDB.CreatedAt,
		UpdatedAt: tenantDB.UpdatedAt,
	}

	// Load the cars information if available
	if tenantDB.Edges.Cars != nil {
		dbModel.Cars = make([]dbmodel.Car, len(tenantDB.Edges.Cars))
		for i, car := range tenantDB.Edges.Cars {
			dbModel.Cars[i] = dbmodel.Car{
				ID:        car.ID,
				TenantID:  car.TenantID,
				Model:     car.Model,
				CreatedAt: car.CreatedAt,
				UpdatedAt: car.UpdatedAt,
			}
		}
		return dbModel.ToDomain(dbmodel.TenantLoadOptions{WithCars: true}), nil
	}

	return dbModel.ToDomain(), nil
}

// Update updates an existing tenant
func (r *tenantRepository) Update(ctx context.Context, tenant *model.Tenant) error {
	_, err := r.client.Tenant.
		UpdateOneID(tenant.ID).
		SetCode(tenant.Code).
		Save(ctx)
	return err
}

// Delete removes a tenant by its ID
func (r *tenantRepository) Delete(ctx context.Context, id string) error {
	return r.client.Tenant.
		DeleteOneID(id).
		Exec(ctx)
}

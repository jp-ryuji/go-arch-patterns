package repository

import (
	"context"

	"github.com/jp-ryuji/go-arch-patterns/internal/domain/entity"
	"github.com/jp-ryuji/go-arch-patterns/internal/domain/repository"
	"github.com/jp-ryuji/go-arch-patterns/internal/infrastructure/postgres/entgen"
	tenant "github.com/jp-ryuji/go-arch-patterns/internal/infrastructure/postgres/entgen/tenant"
)

type tenantRepository struct {
	client *entgen.Client
}

// NewTenantRepository creates a new tenant repository
func NewTenantRepository(client *entgen.Client) repository.TenantRepository {
	return &tenantRepository{client: client}
}

// Create inserts a new tenant into the database
func (r *tenantRepository) Create(ctx context.Context, tenant *entity.Tenant) error {
	_, err := r.client.Tenant.
		Create().
		SetID(tenant.ID).
		SetCode(tenant.Code).
		Save(ctx)
	return err
}

// GetByID retrieves a tenant by its ID
func (r *tenantRepository) GetByID(ctx context.Context, id string) (*entity.Tenant, error) {
	tenantDB, err := r.client.Tenant.
		Query().
		Where(tenant.ID(id)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	// Direct conversion from Ent model to domain entity
	return &entity.Tenant{
		ID:        tenantDB.ID,
		Code:      tenantDB.Code,
		CreatedAt: tenantDB.CreatedAt,
		UpdatedAt: tenantDB.UpdatedAt,
	}, nil
}

// GetByCode retrieves a tenant by its code
func (r *tenantRepository) GetByCode(ctx context.Context, code string) (*entity.Tenant, error) {
	tenantDB, err := r.client.Tenant.
		Query().
		Where(tenant.Code(code)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	// Direct conversion from Ent model to domain entity
	return &entity.Tenant{
		ID:        tenantDB.ID,
		Code:      tenantDB.Code,
		CreatedAt: tenantDB.CreatedAt,
		UpdatedAt: tenantDB.UpdatedAt,
	}, nil
}

// GetByIDWithCars retrieves a tenant by its ID along with its associated cars
func (r *tenantRepository) GetByIDWithCars(ctx context.Context, id string) (*entity.Tenant, error) {
	tenantDB, err := r.client.Tenant.
		Query().
		Where(tenant.ID(id)).
		WithCars().
		Only(ctx)
	if err != nil {
		return nil, err
	}

	// Direct conversion from Ent model to domain entity
	domainTenant := &entity.Tenant{
		ID:        tenantDB.ID,
		Code:      tenantDB.Code,
		CreatedAt: tenantDB.CreatedAt,
		UpdatedAt: tenantDB.UpdatedAt,
	}

	// Load the cars information if available
	if tenantDB.Edges.Cars != nil {
		cars := make(entity.Cars, len(tenantDB.Edges.Cars))
		for i, car := range tenantDB.Edges.Cars {
			cars[i] = &entity.Car{
				ID:        car.ID,
				TenantID:  car.TenantID,
				Model:     car.Model,
				CreatedAt: car.CreatedAt,
				UpdatedAt: car.UpdatedAt,
			}
		}
		domainTenant.Refs = &entity.TenantRefs{
			Cars: cars,
		}
		return domainTenant, nil
	}

	return domainTenant, nil
}

// Update updates an existing tenant
func (r *tenantRepository) Update(ctx context.Context, tenant *entity.Tenant) error {
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

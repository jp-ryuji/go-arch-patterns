package repository

import (
	"context"

	"github.com/jp-ryuji/go-arch-patterns/internal/infrastructure/postgres/entgen"
	tenanttable "github.com/jp-ryuji/go-arch-patterns/internal/infrastructure/postgres/entgen/tenant"
	tenant "github.com/jp-ryuji/go-arch-patterns/internal/tenant"
)

type tenantRepository struct {
	client *entgen.Client
}

func NewTenantRepository(client *entgen.Client) tenant.Repository {
	return &tenantRepository{client: client}
}

func (r *tenantRepository) Create(ctx context.Context, tenantModel *tenant.Tenant) error {
	_, err := r.client.Tenant.
		Create().
		SetID(tenantModel.ID).
		SetCode(tenantModel.Code).
		Save(ctx)
	return err
}

func (r *tenantRepository) GetByID(ctx context.Context, id string) (*tenant.Tenant, error) {
	tenantDB, err := r.client.Tenant.
		Query().
		Where(tenanttable.ID(id)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	return &tenant.Tenant{
		ID:        tenantDB.ID,
		Code:      tenantDB.Code,
		CreatedAt: tenantDB.CreatedAt,
		UpdatedAt: tenantDB.UpdatedAt,
	}, nil
}

func (r *tenantRepository) GetByCode(ctx context.Context, code string) (*tenant.Tenant, error) {
	tenantDB, err := r.client.Tenant.
		Query().
		Where(tenanttable.Code(code)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	return &tenant.Tenant{
		ID:        tenantDB.ID,
		Code:      tenantDB.Code,
		CreatedAt: tenantDB.CreatedAt,
		UpdatedAt: tenantDB.UpdatedAt,
	}, nil
}

func (r *tenantRepository) Update(ctx context.Context, tenantModel *tenant.Tenant) error {
	_, err := r.client.Tenant.
		UpdateOneID(tenantModel.ID).
		SetCode(tenantModel.Code).
		Save(ctx)
	return err
}

func (r *tenantRepository) Delete(ctx context.Context, id string) error {
	return r.client.Tenant.
		DeleteOneID(id).
		Exec(ctx)
}

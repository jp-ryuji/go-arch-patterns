package dbmodel

import (
	"time"

	"github.com/jp-ryuji/go-arch-patterns/internal/domain/entity"
)

// Tenant represents the database model for Tenant
type Tenant struct {
	ID        string
	Code      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	Cars []Car
}

// TenantLoadOptions specifies which associations to load
type TenantLoadOptions struct {
	WithCars bool
}

// ToDomain converts Tenant to domain model with specified associations
func (t *Tenant) ToDomain(opts ...TenantLoadOptions) *entity.Tenant {
	tenant := &entity.Tenant{
		ID:        t.ID,
		Code:      t.Code,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
		Refs:      nil,
	}

	// If no options provided, return basic tenant
	if len(opts) == 0 {
		return tenant
	}

	option := opts[0]

	// Only create Refs if cars need to be loaded
	if option.WithCars && len(t.Cars) > 0 {
		cars := make(entity.Cars, len(t.Cars))
		for i, car := range t.Cars {
			cars[i] = car.ToDomain()
		}

		tenant.Refs = &entity.TenantRefs{
			Cars: cars,
		}
	}

	return tenant
}

// FromDomainTenant converts domain model to Tenant
func FromDomainTenant(tenant *entity.Tenant) *Tenant {
	return &Tenant{
		ID:        tenant.ID,
		Code:      tenant.Code,
		CreatedAt: tenant.CreatedAt,
		UpdatedAt: tenant.UpdatedAt,
	}
}

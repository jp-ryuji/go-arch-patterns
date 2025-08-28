package dbmodel

import (
	"time"

	"github.com/jp-ryuji/go-sample/internal/domain/model"
	"gorm.io/gorm"
)

// Tenant represents the GORM model for Tenant
type Tenant struct {
	ID        string `gorm:"primaryKey;type:varchar(36);not null"`
	Code      string `gorm:"type:varchar(50);not null;uniqueIndex"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Cars []Car `gorm:"foreignKey:TenantID;references:ID"`
}

// TableName specifies the table name for GORM
func (Tenant) TableName() string {
	return "tenants"
}

// TenantLoadOptions specifies which associations to load
type TenantLoadOptions struct {
	WithCars bool
}

// ToDomain converts Tenant to domain model with specified associations
func (t *Tenant) ToDomain(opts ...TenantLoadOptions) *model.Tenant {
	tenant := &model.Tenant{
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
		cars := make(model.Cars, len(t.Cars))
		for i, car := range t.Cars {
			cars[i] = car.ToDomain()
		}

		tenant.Refs = &model.TenantRefs{
			Cars: cars,
		}
	}

	return tenant
}

// FromDomain converts domain model to Tenant
func FromDomainTenant(tenant *model.Tenant) *Tenant {
	return &Tenant{
		ID:        tenant.ID,
		Code:      tenant.Code,
		CreatedAt: tenant.CreatedAt,
		UpdatedAt: tenant.UpdatedAt,
	}
}

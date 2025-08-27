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
}

// TableName specifies the table name for GORM
func (Tenant) TableName() string {
	return "tenants"
}

// ToDomain converts Tenant to domain model
func (t *Tenant) ToDomain() *model.Tenant {
	return &model.Tenant{
		ID:        t.ID,
		Code:      t.Code,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
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

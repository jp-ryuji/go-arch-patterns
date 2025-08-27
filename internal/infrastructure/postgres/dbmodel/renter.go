package dbmodel

import (
	"time"

	"github.com/jp-ryuji/go-sample/internal/domain/model"
	"gorm.io/gorm"
)

// Renter represents the GORM model for Renter
type Renter struct {
	ID                  string `gorm:"primaryKey;type:varchar(36);not null"`
	TenantID            string `gorm:"type:varchar(36);not null;index"`
	ConcreteRenterID    string `gorm:"type:varchar(36);not null;uniqueIndex:idx_renter_unique"`
	ConcreteRenterModel string `gorm:"type:varchar(20);not null;uniqueIndex:idx_renter_unique"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           gorm.DeletedAt `gorm:"index"`
}

// TableName specifies the table name for GORM
func (Renter) TableName() string {
	return "renters"
}

// ToDomain converts Renter to domain model
func (r *Renter) ToDomain() *model.Renter {
	return &model.Renter{
		ID:                  r.ID,
		TenantID:            r.TenantID,
		ConcreteRenterID:    r.ConcreteRenterID,
		ConcreteRenterModel: model.RenterModel(r.ConcreteRenterModel),
		CreatedAt:           r.CreatedAt,
		UpdatedAt:           r.UpdatedAt,
		Refs:                nil, // References would be loaded separately if needed
	}
}

// FromDomain converts domain model to Renter
func FromDomainRenter(renter *model.Renter) *Renter {
	return &Renter{
		ID:                  renter.ID,
		TenantID:            renter.TenantID,
		ConcreteRenterID:    renter.ConcreteRenterID,
		ConcreteRenterModel: string(renter.ConcreteRenterModel),
		CreatedAt:           renter.CreatedAt,
		UpdatedAt:           renter.UpdatedAt,
	}
}

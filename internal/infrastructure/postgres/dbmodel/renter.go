package dbmodel

import (
	"time"

	"github.com/jp-ryuji/go-sample/internal/domain/model"
	"gorm.io/gorm"
)

// Renter represents the GORM model for Renter
type Renter struct {
	ID               string `gorm:"primaryKey;type:varchar(36);not null"`
	TenantID         string `gorm:"type:varchar(36);not null;index"`
	RenterEntityID   string `gorm:"type:varchar(36);not null"`
	RenterEntityType string `gorm:"type:varchar(20);not null"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`

	Rentals []Rental `gorm:"foreignKey:RenterID"`
}

// TableName specifies the table name for GORM
func (Renter) TableName() string {
	return "renters"
}

// RenterLoadOptions specifies which associations to load
type RenterLoadOptions struct {
	WithCompany    bool
	WithIndividual bool
	WithRentals    bool
}

// ToDomain converts Renter to domain model with specified associations
func (r *Renter) ToDomain(opts ...RenterLoadOptions) *model.Renter {
	renter := &model.Renter{
		ID:               r.ID,
		TenantID:         r.TenantID,
		RenterEntityID:   r.RenterEntityID,
		RenterEntityType: model.RenterEntity(r.RenterEntityType),
		CreatedAt:        r.CreatedAt,
		UpdatedAt:        r.UpdatedAt,
		Refs:             nil,
	}

	// If no options provided, return basic renter
	if len(opts) == 0 {
		return renter
	}

	option := opts[0]

	// Create Refs if any associations need to be loaded
	if option.WithCompany || option.WithIndividual || option.WithRentals {
		renter.Refs = &model.RenterRefs{}

		// Load company if requested and available
		// Note: This requires the Company association to be preloaded
		// We're not directly accessing a Company field here because GORM doesn't support
		// polymorphic associations in the same way. This would need to be handled at a higher level.
		_ = option.WithCompany && r.RenterEntityType == "company" && r.RenterEntityID != ""

		// Load individual if requested and available
		// Note: This requires the Individual association to be preloaded
		// We're not directly accessing an Individual field here because GORM doesn't support
		// polymorphic associations in the same way. This would need to be handled at a higher level.
		_ = option.WithIndividual && r.RenterEntityType == "individual" && r.RenterEntityID != ""

		// Load rentals if requested and available
		if option.WithRentals && len(r.Rentals) > 0 {
			rentals := make(model.Rentals, len(r.Rentals))
			for i, rental := range r.Rentals {
				rentals[i] = rental.ToDomain()
			}
			renter.Refs.Rentals = rentals
		}
	}

	return renter
}

// FromDomain converts domain model to Renter
func FromDomainRenter(renter *model.Renter) *Renter {
	return &Renter{
		ID:               renter.ID,
		TenantID:         renter.TenantID,
		RenterEntityID:   renter.RenterEntityID,
		RenterEntityType: string(renter.RenterEntityType),
		CreatedAt:        renter.CreatedAt,
		UpdatedAt:        renter.UpdatedAt,
	}
}

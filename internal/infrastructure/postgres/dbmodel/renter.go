package dbmodel

import (
	"time"

	"github.com/jp-ryuji/go-sample/internal/domain/model"
)

// Renter represents the database model for Renter
type Renter struct {
	ID               string     `json:"id"`
	TenantID         string     `json:"tenant_id"`
	RenterEntityID   string     `json:"renter_entity_id"`
	RenterEntityType string     `json:"renter_entity_type"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty"`

	Rentals []Rental `json:"rentals,omitempty"`
}

// TableName specifies the table name (kept for compatibility with existing code)
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
		// Polymorphic associations require special handling at a higher level.
		_ = option.WithCompany && r.RenterEntityType == "company" && r.RenterEntityID != ""

		// Load individual if requested and available
		// Note: This requires the Individual association to be preloaded
		// Polymorphic associations require special handling at a higher level.
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

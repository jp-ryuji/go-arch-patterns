package dbmodel

import (
	"time"

	"github.com/jp-ryuji/go-arch-patterns/internal/domain/model"
)

// Renter represents the database model for Renter
type Renter struct {
	ID        string
	TenantID  string
	Type      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	Rentals []Rental
}

// RenterLoadOptions specifies which associations to load
type RenterLoadOptions struct {
	WithRentals bool
}

// ToDomain converts Renter to domain model with specified associations
func (r *Renter) ToDomain(opts ...RenterLoadOptions) *model.Renter {
	renter := &model.Renter{
		ID:        r.ID,
		TenantID:  r.TenantID,
		Type:      model.RenterType(r.Type),
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
		Refs:      nil,
	}

	// If no options provided, return basic renter
	if len(opts) == 0 {
		return renter
	}

	option := opts[0]

	// Create Refs if any associations need to be loaded
	if option.WithRentals {
		renter.Refs = &model.RenterRefs{}

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
		ID:        renter.ID,
		TenantID:  renter.TenantID,
		Type:      string(renter.Type),
		CreatedAt: renter.CreatedAt,
		UpdatedAt: renter.UpdatedAt,
	}
}

package dbmodel

import (
	"time"

	"github.com/jp-ryuji/go-sample/internal/domain/model"
)

// Rental represents the database model for Rental
type Rental struct {
	ID        string     `json:"id"`
	TenantID  string     `json:"tenant_id"`
	CarID     string     `json:"car_id"`
	RenterID  string     `json:"renter_id"`
	StartsAt  time.Time  `json:"starts_at"`
	EndsAt    time.Time  `json:"ends_at"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	Tenant        Tenant         `json:"tenant,omitempty"`
	Car           Car            `json:"car,omitempty"`
	Renter        Renter         `json:"renter,omitempty"`
	RentalOptions []RentalOption `json:"rental_options,omitempty"`
}

// TableName specifies the table name (kept for compatibility with existing code)
func (Rental) TableName() string {
	return "rentals"
}

// RentalLoadOptions specifies which associations to load
type RentalLoadOptions struct {
	WithTenant        bool
	WithCar           bool
	WithRenter        bool
	WithRentalOptions bool
}

// ToDomain converts Rental to domain model with specified associations
func (r *Rental) ToDomain(opts ...RentalLoadOptions) *model.Rental {
	rental := &model.Rental{
		ID:        r.ID,
		TenantID:  r.TenantID,
		CarID:     r.CarID,
		RenterID:  r.RenterID,
		StartsAt:  r.StartsAt,
		EndsAt:    r.EndsAt,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
		Refs:      nil,
	}

	// If no options provided, return basic rental
	if len(opts) == 0 {
		return rental
	}

	option := opts[0]

	// Only create Refs if at least one association needs to be loaded
	if option.WithTenant || option.WithCar || option.WithRenter || option.WithRentalOptions {
		rental.Refs = &model.RentalRefs{}

		if option.WithTenant && r.Tenant.ID != "" {
			tenant := r.Tenant.ToDomain()
			rental.Refs.Tenant = tenant
		}

		if option.WithCar && r.Car.ID != "" {
			car := r.Car.ToDomain()
			rental.Refs.Car = car
		}

		if option.WithRenter && r.Renter.ID != "" {
			// When loading a renter, we need to determine if we should also load its associated entity
			renterOpts := RenterLoadOptions{}
			// For now, we're not automatically loading the Company/Individual associations
			// This would need to be handled at a higher level or with additional parameters
			renter := r.Renter.ToDomain(renterOpts)
			rental.Refs.Renter = renter
		}

		if option.WithRentalOptions && len(r.RentalOptions) > 0 {
			rentalOptions := make(model.RentalOptions, len(r.RentalOptions))
			for i, rentalOption := range r.RentalOptions {
				rentalOptions[i] = rentalOption.ToDomain()
			}
			rental.Refs.RentalOptions = rentalOptions
		}
	}

	return rental
}

// FromDomain converts domain model to Rental
func FromDomainRental(rental *model.Rental) *Rental {
	return &Rental{
		ID:        rental.ID,
		TenantID:  rental.TenantID,
		CarID:     rental.CarID,
		RenterID:  rental.RenterID,
		StartsAt:  rental.StartsAt,
		EndsAt:    rental.EndsAt,
		CreatedAt: rental.CreatedAt,
		UpdatedAt: rental.UpdatedAt,
	}
}

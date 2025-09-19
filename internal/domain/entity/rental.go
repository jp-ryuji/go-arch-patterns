package entity

import (
	"time"

	"github.com/oklog/ulid/v2"
)

// Rentals is a slice of Rental
type Rentals []*Rental

// Rental represents a rental entity
type Rental struct {
	ID        string
	TenantID  string
	CarID     string
	RenterID  string
	StartsAt  time.Time
	EndsAt    time.Time
	CreatedAt time.Time
	UpdatedAt time.Time

	// References to related entities
	Refs *RentalRefs
}

// RentalRefs holds references to related entities
type RentalRefs struct {
	Tenant        *Tenant
	Car           *Car
	Renter        *Renter
	RentalOptions RentalOptions
}

// NewRental creates a new Rental
func NewRental(tenantID, carID, renterID string, startsAt, endsAt time.Time) *Rental {
	now := time.Now()
	return &Rental{
		ID:        ulid.Make().String(),
		TenantID:  tenantID,
		CarID:     carID,
		RenterID:  renterID,
		StartsAt:  startsAt,
		EndsAt:    endsAt,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// WithID creates a Rental with a specific ID (for testing)
func (r *Rental) WithID(id string) *Rental {
	r.ID = id
	return r
}

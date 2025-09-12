package model

import (
	"time"

	"github.com/oklog/ulid/v2"
)

// Renters is a slice of Renter
type Renters []*Renter

// RenterType represents the type of renter
type RenterType string

const (
	CompanyRenter    RenterType = "company"
	IndividualRenter RenterType = "individual"
)

// Renter represents a renter entity (base class for Company and Individual)
type Renter struct {
	ID        string
	TenantID  string
	Type      RenterType
	CreatedAt time.Time
	UpdatedAt time.Time

	// References to related entities
	Refs *RenterRefs
}

// RenterRefs holds references to related entities
type RenterRefs struct {
	Rentals Rentals
}

// NewRenter creates a new Renter
func NewRenter(tenantID string, renterType RenterType, createdAt time.Time) *Renter {
	return &Renter{
		ID:        ulid.Make().String(),
		TenantID:  tenantID,
		Type:      renterType,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}
}

// WithID creates a Renter with a specific ID (for testing)
func (r *Renter) WithID(id string) *Renter {
	r.ID = id
	return r
}

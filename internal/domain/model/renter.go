package model

import (
	"time"

	"github.com/google/uuid"
)

// Renters is a slice of Renter
type Renters []*Renter

// RenterEntity represents the type of entity that can be a renter
type RenterEntity string

const (
	CompanyRenter    RenterEntity = "company"
	IndividualRenter RenterEntity = "individual"
)

// Renter represents a renter entity
type Renter struct {
	ID               string
	TenantID         string
	RenterEntityID   string
	RenterEntityType RenterEntity
	CreatedAt        time.Time
	UpdatedAt        time.Time

	// References to related entities
	Refs *RenterRefs
}

// RenterRefs holds references to related entities
type RenterRefs struct {
	Company    *Company
	Individual *Individual
	Rentals    Rentals
}

// NewRenter creates a new Renter
func NewRenter(tenantID, renterEntityID string, renterEntityType RenterEntity, createdAt time.Time) *Renter {
	return &Renter{
		ID:               uuid.New().String(),
		TenantID:         tenantID,
		RenterEntityID:   renterEntityID,
		RenterEntityType: renterEntityType,
		CreatedAt:        createdAt,
		UpdatedAt:        createdAt,
	}
}

// WithID creates a Renter with a specific ID (for testing)
func (r *Renter) WithID(id string) *Renter {
	r.ID = id
	return r
}

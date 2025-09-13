package model

import (
	"time"

	"github.com/aarondl/null/v9"
	"github.com/oklog/ulid/v2"

	"github.com/jp-ryuji/go-arch-patterns/internal/domain/model/value"
)

// Individuals is a slice of Individual
type Individuals []*Individual

// Individual represents an individual entity (inherits from Renter)
type Individual struct {
	ID        string
	RenterID  string
	TenantID  string
	Email     value.Email
	FirstName null.String
	LastName  null.String
	CreatedAt time.Time
	UpdatedAt time.Time

	// References to related entities
	Refs *IndividualRefs
}

// IndividualRefs holds references to related entities
type IndividualRefs struct{}

// NewIndividual creates a new Individual
func NewIndividual(renterID, tenantID string, email value.Email, firstName, lastName null.String, createdAt time.Time) *Individual {
	return &Individual{
		ID:        ulid.Make().String(),
		RenterID:  renterID,
		TenantID:  tenantID,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}
}

// WithRenterID creates an Individual with a specific RenterID (for testing)
func (i *Individual) WithRenterID(renterID string) *Individual {
	i.RenterID = renterID
	return i
}

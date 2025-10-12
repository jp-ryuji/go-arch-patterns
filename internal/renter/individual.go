package renter

import (
	"time"

	"github.com/aarondl/null/v9"
	"github.com/oklog/ulid/v2"
)

type Individuals []*Individual

type Individual struct {
	ID        string
	RenterID  string
	TenantID  string
	Email     Email
	FirstName null.String
	LastName  null.String
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewIndividual(renterID, tenantID string, email Email, firstName, lastName null.String, createdAt time.Time) *Individual {
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

func (i *Individual) WithRenterID(renterID string) *Individual {
	i.RenterID = renterID
	return i
}

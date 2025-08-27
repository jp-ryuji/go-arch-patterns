package model

import (
	"time"

	"github.com/jp-ryuji/go-sample/internal/pkg/id"
)

// RentalableRefs holds references to related entities for Rentalable
type RentalableRefs struct {
	Tenant *Tenant
	// Renter could be either a Company or an Individual
	Renter Renter
}

type Rentalable struct {
	ID          string
	TenantID    string
	RenterID    string
	RenterModel RenterModel
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Refs *RentalableRefs
}

type Rentalables []*Rentalable

func NewRentalable(
	tenantID string,
	renterID string,
	renterModel RenterModel,
	t time.Time,
) *Rentalable {
	return &Rentalable{
		ID:          id.New(),
		TenantID:    tenantID,
		RenterID:    renterID,
		RenterModel: renterModel,
		CreatedAt:   t,
		UpdatedAt:   t,

		Refs: nil,
	}
}

type RenterModel string

const (
	RenterModelCompany    RenterModel = "company"
	RenterModelIndividual RenterModel = "individual"
)

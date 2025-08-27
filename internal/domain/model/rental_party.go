package model

import (
	"time"

	"github.com/jp-ryuji/go-sample/internal/pkg/id"
)

// RentalPartyRefs holds references to related entities for RentalParty
type RentalPartyRefs struct {
	Tenant *Tenant
	// Renter could be either a Company or an Individual
	Renter Renter
	// Rentals associated with this RentalParty
	Rentals Rentals
}

type RentalParty struct {
	ID          string
	TenantID    string
	RenterID    string
	RenterModel RenterModel
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Refs *RentalPartyRefs
}

type RentalParties []*RentalParty

func NewRentalParty(
	tenantID string,
	renterID string,
	renterModel RenterModel,
	t time.Time,
) *RentalParty {
	return &RentalParty{
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

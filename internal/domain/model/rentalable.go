package model

import (
	"time"

	"github.com/jp-ryuji/go-sample/internal/pkg/id"
)

// RentalableRefs holds references to related entities for Rentalable
type RentalableRefs struct {
	Tenant *Tenant
	// Renter could be either a Company or an Individual
	// This is a common pattern in DDD for polymorphic associations
	// The actual type would need to be determined based on RenterModel
	// Consider using interfaces or a sum type pattern if Go allows for cleaner implementation
	Renter interface{} // This is a placeholder, consider a better approach
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

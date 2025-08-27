package model

import (
	"time"

	"github.com/jp-ryuji/go-sample/internal/pkg/id"
)

// RenterableRefs holds references to related entities for Renterable
type RenterableRefs struct {
	Tenant *Tenant
	// Renter could be either a Company or an Individual
	Renter Renter
}

type Renterable struct {
	ID          string
	TenantID    string
	RenterID    string
	RenterModel RenterModel
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Refs *RenterableRefs
}

type Renterables []*Renterable

func NewRenterable(
	tenantID string,
	renterID string,
	renterModel RenterModel,
	t time.Time,
) *Renterable {
	return &Renterable{
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

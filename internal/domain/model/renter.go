package model

import (
	"time"

	"github.com/jp-ryuji/go-sample/internal/pkg/id"
)

// ConcreteRenter is an interface that represents a concrete renter entity.
// Both Company and Individual implement this interface.
type ConcreteRenter interface {
	// GetID returns the unique identifier of the concrete renter.
	GetID() string
	// GetTenantID returns the tenant identifier associated with the concrete renter.
	GetTenantID() string
	// GetType returns the type of the concrete renter (e.g., "company" or "individual").
	GetType() RenterModel
}

// RenterRefs holds references to related entities for the Renter association.
type RenterRefs struct {
	Tenant *Tenant
	// ConcreteRenter is the actual Company or Individual fulfilling the renter role.
	ConcreteRenter ConcreteRenter
	// Rentals associated with this Renter link.
	Rentals Rentals
}

// Renter represents the association between a Rental and a specific Renter (Company or Individual).
// It holds the information needed to link to the concrete renter entity.
type Renter struct {
	// ID is the unique identifier for this Renter association record.
	ID       string
	TenantID string
	// ConcreteRenterID is the ID of the actual Company or Individual.
	ConcreteRenterID string
	// ConcreteRenterModel specifies the type ("company" or "individual") of the ConcreteRenter.
	ConcreteRenterModel RenterModel
	CreatedAt           time.Time
	UpdatedAt           time.Time

	// Refs holds lazily/eagerly loaded references. It should not be populated directly.
	Refs *RenterRefs
}

// Renters is a slice of pointers to Renter.
type Renters []*Renter

// NewRenter creates a new instance of Renter.
func NewRenter(
	tenantID string,
	concreteRenterID string,
	concreteRenterModel RenterModel,
	t time.Time,
) *Renter {
	return &Renter{
		ID:                  id.New(),
		TenantID:            tenantID,
		ConcreteRenterID:    concreteRenterID,
		ConcreteRenterModel: concreteRenterModel,
		CreatedAt:           t,
		UpdatedAt:           t,
		Refs:                nil,
	}
}

// RenterModel defines the type of the concrete renter.
type RenterModel string

const (
	// RenterModelCompany indicates the renter is a Company.
	RenterModelCompany RenterModel = "company"
	// RenterModelIndividual indicates the renter is an Individual.
	RenterModelIndividual RenterModel = "individual"
)

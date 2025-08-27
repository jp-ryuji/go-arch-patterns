package model

import (
	"time"

	"github.com/jp-ryuji/go-sample/internal/pkg/id"
)

// RenterEntityInstance is an interface for renter entities.
type RenterEntityInstance interface {
	GetID() string
	GetEntityType() RenterEntity
}

// RenterRefs holds references to related entities for the Renter association.
type RenterRefs struct {
	Tenant               *Tenant
	RenterEntityInstance RenterEntityInstance // the actual Company or Individual fulfilling the renter role
	Rentals              Rentals              // rentals associated with this Renter link
}

// Renter represents the association between a Rental and a specific Renter.
type Renter struct {
	ID               string
	TenantID         string
	RenterEntityID   string       // ID of the actual Company or Individual
	RenterEntityType RenterEntity // type ("company" or "individual") of the entity
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Refs             *RenterRefs // lazily/eagerly loaded references
}

// Renters is a slice of pointers to Renter.
type Renters []*Renter

// NewRenter creates a new instance of Renter.
func NewRenter(
	tenantID string,
	renterEntityID string,
	renterEntityType RenterEntity,
	t time.Time,
) *Renter {
	return &Renter{
		ID:               id.New(),
		TenantID:         tenantID,
		RenterEntityID:   renterEntityID,
		RenterEntityType: renterEntityType,
		CreatedAt:        t,
		UpdatedAt:        t,
		Refs:             nil,
	}
}

// RenterEntity defines the type of the renter entity.
type RenterEntity string

const (
	RenterEntityCompany    RenterEntity = "company"
	RenterEntityIndividual RenterEntity = "individual"
)

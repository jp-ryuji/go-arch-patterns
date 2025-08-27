package model

import (
	"time"

	"github.com/jp-ryuji/go-sample/internal/pkg/id"
)

// RentalRefs holds references to related entities for Rental
type RentalRefs struct {
	Tenant      *Tenant
	Car         *Car
	RentalParty *RentalParty
}

type Rental struct {
	ID            string
	TenantID      string
	CarID         string
	RentalPartyID string
	StartsAt      time.Time
	EndsAt        time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time

	Refs *RentalRefs
}

type Rentals []*Rental

func NewRental(
	tenantID string,
	rentalPartyID string,
	startsAt time.Time,
	endsAt time.Time,
	t time.Time,
) *Rental {
	return &Rental{
		ID:            id.New(),
		TenantID:      tenantID,
		RentalPartyID: rentalPartyID,
		StartsAt:      startsAt,
		EndsAt:        endsAt,
		CreatedAt:     t,
		UpdatedAt:     t,

		Refs: nil,
	}
}

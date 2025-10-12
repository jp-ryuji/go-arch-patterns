package rental

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type Rentals []*Rental

type Rental struct {
	ID        string
	TenantID  string
	CarID     string
	RenterID  string
	StartsAt  time.Time
	EndsAt    time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewRental(tenantID, carID, renterID string, startsAt, endsAt time.Time) *Rental {
	now := time.Now()
	return &Rental{
		ID:        ulid.Make().String(),
		TenantID:  tenantID,
		CarID:     carID,
		RenterID:  renterID,
		StartsAt:  startsAt,
		EndsAt:    endsAt,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (r *Rental) WithID(id string) *Rental {
	r.ID = id
	return r
}

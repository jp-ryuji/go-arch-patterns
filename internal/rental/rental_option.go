package rental

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type RentalOptions []*RentalOption

type RentalOption struct {
	ID        string
	TenantID  string
	RentalID  string
	OptionID  string
	Count     int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewRentalOption(tenantID, rentalID, optionID string, count int) *RentalOption {
	now := time.Now()
	return &RentalOption{
		ID:        ulid.Make().String(),
		TenantID:  tenantID,
		RentalID:  rentalID,
		OptionID:  optionID,
		Count:     count,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (ro *RentalOption) WithID(id string) *RentalOption {
	ro.ID = id
	return ro
}

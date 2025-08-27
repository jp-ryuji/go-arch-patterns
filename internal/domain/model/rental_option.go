package model

import (
	"time"

	"github.com/jp-ryuji/go-sample/internal/pkg/id"
)

// RentalOptionRefs holds references to related entities for RentalOption
type RentalOptionRefs struct {
	Tenant *Tenant
	Rental *Rental
	Option *Option
}

type RentalOption struct {
	ID        string
	TenantID  string
	RentalID  string
	OptionID  string
	Count     int
	CreatedAt time.Time
	UpdatedAt time.Time

	Refs *RentalOptionRefs
}

type RentalOptions []*RentalOption

func NewRentalOption(
	tenantID string,
	rentalID string,
	optionID string,
	count int,
	t time.Time,
) *RentalOption {
	return &RentalOption{
		ID:        id.New(),
		TenantID:  tenantID,
		RentalID:  rentalID,
		OptionID:  optionID,
		Count:     count,
		CreatedAt: t,
		UpdatedAt: t,

		Refs: nil,
	}
}

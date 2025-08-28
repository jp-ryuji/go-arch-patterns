package model

import (
	"time"

	"github.com/google/uuid"
)

// RentalOptions is a slice of RentalOption
type RentalOptions []*RentalOption

// RentalOption represents a rental option entity
type RentalOption struct {
	ID        string
	TenantID  string
	RentalID  string
	OptionID  string
	Count     int
	CreatedAt time.Time
	UpdatedAt time.Time

	// References to related entities
	Refs *RentalOptionRefs
}

// RentalOptionRefs holds references to related entities
type RentalOptionRefs struct {
	Rental *Rental
	Option *Option
}

// NewRentalOption creates a new RentalOption
func NewRentalOption(tenantID, rentalID, optionID string, count int) *RentalOption {
	now := time.Now()
	return &RentalOption{
		ID:        uuid.New().String(),
		TenantID:  tenantID,
		RentalID:  rentalID,
		OptionID:  optionID,
		Count:     count,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// WithID creates a RentalOption with a specific ID (for testing)
func (ro *RentalOption) WithID(id string) *RentalOption {
	ro.ID = id
	return ro
}

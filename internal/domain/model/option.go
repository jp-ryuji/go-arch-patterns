package model

import (
	"time"

	"github.com/oklog/ulid/v2"
)

// Options is a slice of Option
type Options []*Option

// Option represents an option entity
type Option struct {
	ID        string
	TenantID  string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time

	// References to related entities
	Refs *OptionRefs
}

// OptionRefs holds references to related entities
type OptionRefs struct {
	RentalOptions RentalOptions
}

// NewOption creates a new Option
func NewOption(tenantID, name string) *Option {
	now := time.Now()
	return &Option{
		ID:        ulid.Make().String(),
		TenantID:  tenantID,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// WithID creates an Option with a specific ID (for testing)
func (o *Option) WithID(id string) *Option {
	o.ID = id
	return o
}

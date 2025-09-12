package model

import (
	"time"

	"github.com/oklog/ulid/v2"
)

// Tenants is a slice of Tenant
type Tenants []*Tenant

// Tenant represents a tenant entity
type Tenant struct {
	ID        string
	Code      string
	CreatedAt time.Time
	UpdatedAt time.Time

	// References to related entities
	Refs *TenantRefs
}

// TenantRefs holds references to related entities
type TenantRefs struct {
	Cars Cars
}

// NewTenant creates a new Tenant
func NewTenant(code string, createdAt time.Time) *Tenant {
	return &Tenant{
		ID:        ulid.Make().String(),
		Code:      code,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}
}

// WithID creates a Tenant with a specific ID (for testing)
func (t *Tenant) WithID(id string) *Tenant {
	t.ID = id
	return t
}

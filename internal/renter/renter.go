package renter

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type Type string

type Renters []*Renter

const (
	TypeCompany    Type = "company"
	TypeIndividual Type = "individual"
)

type Renter struct {
	ID        string
	TenantID  string
	Type      Type
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewRenter(tenantID string, renterType Type, createdAt time.Time) *Renter {
	return &Renter{
		ID:        ulid.Make().String(),
		TenantID:  tenantID,
		Type:      renterType,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}
}

func (r *Renter) WithID(id string) *Renter {
	r.ID = id
	return r
}

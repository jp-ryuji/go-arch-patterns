package model

import (
	"time"

	"github.com/jp-ryuji/go-sample/internal/pkg/id"
)

type Tenant struct {
	ID        string
	Code      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Tenants []*Tenant

func NewTenant(
	code string,
	t time.Time,
) *Tenant {
	return &Tenant{
		ID:        id.New(),
		Code:      code,
		CreatedAt: t,
		UpdatedAt: t,
	}
}

package model

import (
	"time"

	"github.com/jp-ryuji/go-sample/internal/pkg/id"
)

// CarRefs holds references to related entities for Car
type CarRefs struct {
	Tenant *Tenant
}

type Car struct {
	ID        string
	TenantID  string
	Model     string
	CreatedAt time.Time
	UpdatedAt time.Time

	Refs *CarRefs
}

type Cars []*Car

func NewCar(
	tenantID string,
	model string,
	t time.Time,
) *Car {
	return &Car{
		ID:        id.New(),
		TenantID:  tenantID,
		Model:     model,
		CreatedAt: t,
		UpdatedAt: t,

		Refs: nil,
	}
}

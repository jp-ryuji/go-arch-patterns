package rental

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type Cars []*Car

type Car struct {
	ID        string
	TenantID  string
	Model     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewCar(tenantID, model string, createdAt time.Time) *Car {
	return &Car{
		ID:        ulid.Make().String(),
		TenantID:  tenantID,
		Model:     model,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}
}

func (c *Car) WithID(id string) *Car {
	c.ID = id
	return c
}

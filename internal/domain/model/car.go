package model

import (
	"time"

	"github.com/google/uuid"
)

// Cars is a slice of Car
type Cars []*Car

// Car represents a car entity
type Car struct {
	ID        string
	TenantID  string
	Model     string
	CreatedAt time.Time
	UpdatedAt time.Time

	// References to related entities
	Refs *CarRefs
}

// CarRefs holds references to related entities
type CarRefs struct {
	Tenant  *Tenant
	Rentals Rentals
}

// NewCar creates a new Car
func NewCar(tenantID, model string, createdAt time.Time) *Car {
	return &Car{
		ID:        uuid.New().String(),
		TenantID:  tenantID,
		Model:     model,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}
}

// WithID creates a Car with a specific ID (for testing)
func (c *Car) WithID(id string) *Car {
	c.ID = id
	return c
}

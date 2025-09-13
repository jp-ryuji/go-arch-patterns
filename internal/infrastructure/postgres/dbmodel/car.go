package dbmodel

import (
	"time"

	"github.com/jp-ryuji/go-arch-patterns/internal/domain/model"
)

// Car represents the database model for Car
type Car struct {
	ID        string
	TenantID  string
	Model     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	Tenant  Tenant
	Rentals []Rental
}

// CarLoadOptions specifies which associations to load
type CarLoadOptions struct {
	WithTenant  bool
	WithRentals bool
}

// ToDomain converts Car to domain model with specified associations
func (c *Car) ToDomain(opts ...CarLoadOptions) *model.Car {
	car := &model.Car{
		ID:        c.ID,
		TenantID:  c.TenantID,
		Model:     c.Model,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Refs:      nil,
	}

	// If no options provided, return basic car
	if len(opts) == 0 {
		return car
	}

	option := opts[0]

	// Only create Refs if tenant needs to be loaded and is available, or rentals need to be loaded and are available
	if (option.WithTenant && c.Tenant.ID != "") || (option.WithRentals && len(c.Rentals) > 0) {
		car.Refs = &model.CarRefs{}

		// Load tenant if requested and available
		if option.WithTenant && c.Tenant.ID != "" {
			car.Refs.Tenant = c.Tenant.ToDomain()
		}

		// Load rentals if requested and available
		if option.WithRentals && len(c.Rentals) > 0 {
			rentals := make(model.Rentals, len(c.Rentals))
			for i, rental := range c.Rentals {
				rentals[i] = rental.ToDomain()
			}
			car.Refs.Rentals = rentals
		}
	}

	return car
}

// FromDomain converts domain model to Car
func FromDomainCar(car *model.Car) *Car {
	return &Car{
		ID:        car.ID,
		TenantID:  car.TenantID,
		Model:     car.Model,
		CreatedAt: car.CreatedAt,
		UpdatedAt: car.UpdatedAt,
	}
}

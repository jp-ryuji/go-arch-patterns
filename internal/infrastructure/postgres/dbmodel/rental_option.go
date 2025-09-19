package dbmodel

import (
	"time"

	"github.com/jp-ryuji/go-arch-patterns/internal/domain/entity"
)

// RentalOption represents the database model for RentalOption
type RentalOption struct {
	ID        string
	TenantID  string
	RentalID  string
	OptionID  string
	Count     int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	Rental Rental
	Option Option
}

// RentalOptionLoadOptions specifies which associations to load
type RentalOptionLoadOptions struct {
	WithRental bool
	WithOption bool
}

// ToDomain converts RentalOption to domain model with specified associations
func (ro *RentalOption) ToDomain(opts ...RentalOptionLoadOptions) *entity.RentalOption {
	rentalOption := &entity.RentalOption{
		ID:        ro.ID,
		TenantID:  ro.TenantID,
		RentalID:  ro.RentalID,
		OptionID:  ro.OptionID,
		Count:     ro.Count,
		CreatedAt: ro.CreatedAt,
		UpdatedAt: ro.UpdatedAt,
		Refs:      nil,
	}

	// If no options provided, return basic rental option
	if len(opts) == 0 {
		return rentalOption
	}

	option := opts[0]

	// Create Refs if any associations need to be loaded
	if option.WithRental || option.WithOption {
		rentalOption.Refs = &entity.RentalOptionRefs{}

		// Load rental if requested and available
		if option.WithRental && ro.Rental.ID != "" {
			rentalOption.Refs.Rental = ro.Rental.ToDomain()
		}

		// Load option if requested and available
		if option.WithOption && ro.Option.ID != "" {
			rentalOption.Refs.Option = ro.Option.ToDomain()
		}
	}

	return rentalOption
}

// FromDomainRentalOption converts domain model to RentalOption
func FromDomainRentalOption(rentalOption *entity.RentalOption) *RentalOption {
	return &RentalOption{
		ID:        rentalOption.ID,
		TenantID:  rentalOption.TenantID,
		RentalID:  rentalOption.RentalID,
		OptionID:  rentalOption.OptionID,
		Count:     rentalOption.Count,
		CreatedAt: rentalOption.CreatedAt,
		UpdatedAt: rentalOption.UpdatedAt,
	}
}

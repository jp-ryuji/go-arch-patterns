package dbmodel

import (
	"time"

	"github.com/jp-ryuji/go-sample/internal/domain/model"
	"gorm.io/gorm"
)

// RentalOption represents the GORM model for RentalOption
type RentalOption struct {
	ID        string `gorm:"primaryKey;type:varchar(36);not null"`
	TenantID  string `gorm:"type:varchar(36);not null;index"`
	RentalID  string `gorm:"type:varchar(36);not null;uniqueIndex:idx_rental_option_unique"`
	OptionID  string `gorm:"type:varchar(36);not null;uniqueIndex:idx_rental_option_unique"`
	Count     int    `gorm:"type:integer;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Rental Rental `gorm:"foreignKey:RentalID"`
	Option Option `gorm:"foreignKey:OptionID"`
}

// TableName specifies the table name for GORM
func (RentalOption) TableName() string {
	return "rental_options"
}

// RentalOptionLoadOptions specifies which associations to load
type RentalOptionLoadOptions struct {
	WithRental bool
	WithOption bool
}

// ToDomain converts RentalOption to domain model with specified associations
func (ro *RentalOption) ToDomain(opts ...RentalOptionLoadOptions) *model.RentalOption {
	rentalOption := &model.RentalOption{
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
		rentalOption.Refs = &model.RentalOptionRefs{}

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

// FromDomain converts domain model to RentalOption
func FromDomainRentalOption(rentalOption *model.RentalOption) *RentalOption {
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

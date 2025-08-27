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
}

// TableName specifies the table name for GORM
func (RentalOption) TableName() string {
	return "rental_options"
}

// ToDomain converts RentalOption to domain model
func (ro *RentalOption) ToDomain() *model.RentalOption {
	return &model.RentalOption{
		ID:        ro.ID,
		TenantID:  ro.TenantID,
		RentalID:  ro.RentalID,
		OptionID:  ro.OptionID,
		Count:     ro.Count,
		CreatedAt: ro.CreatedAt,
		UpdatedAt: ro.UpdatedAt,
		Refs:      nil, // References would be loaded separately if needed
	}
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

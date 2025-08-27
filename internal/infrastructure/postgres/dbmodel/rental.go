package dbmodel

import (
	"time"

	"github.com/jp-ryuji/go-sample/internal/domain/model"
	"gorm.io/gorm"
)

// Rental represents the GORM model for Rental
type Rental struct {
	ID        string `gorm:"primaryKey;type:varchar(36);not null"`
	TenantID  string `gorm:"type:varchar(36);not null;index"`
	CarID     string `gorm:"type:varchar(36);not null;index"`
	RenterID  string `gorm:"type:varchar(36);not null;index"`
	StartsAt  time.Time
	EndsAt    time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName specifies the table name for GORM
func (Rental) TableName() string {
	return "rentals"
}

// ToDomain converts Rental to domain model
func (r *Rental) ToDomain() *model.Rental {
	return &model.Rental{
		ID:        r.ID,
		TenantID:  r.TenantID,
		CarID:     r.CarID,
		RenterID:  r.RenterID,
		StartsAt:  r.StartsAt,
		EndsAt:    r.EndsAt,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
		Refs:      nil, // References would be loaded separately if needed
	}
}

// FromDomain converts domain model to Rental
func FromDomainRental(rental *model.Rental) *Rental {
	return &Rental{
		ID:        rental.ID,
		TenantID:  rental.TenantID,
		CarID:     rental.CarID,
		RenterID:  rental.RenterID,
		StartsAt:  rental.StartsAt,
		EndsAt:    rental.EndsAt,
		CreatedAt: rental.CreatedAt,
		UpdatedAt: rental.UpdatedAt,
	}
}

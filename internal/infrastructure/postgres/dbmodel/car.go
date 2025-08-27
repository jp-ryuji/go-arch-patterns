package dbmodel

import (
	"time"

	"github.com/jp-ryuji/go-sample/internal/domain/model"
	"gorm.io/gorm"
)

// Car represents the GORM model for Car
type Car struct {
	ID        string `gorm:"primaryKey;type:varchar(36);not null"`
	TenantID  string `gorm:"type:varchar(36);not null;uniqueIndex:idx_car_tenant_model"`
	Model     string `gorm:"type:varchar(255);not null;uniqueIndex:idx_car_tenant_model"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName specifies the table name for GORM
func (Car) TableName() string {
	return "cars"
}

// ToDomain converts Car to domain model
func (c *Car) ToDomain() *model.Car {
	return &model.Car{
		ID:        c.ID,
		TenantID:  c.TenantID,
		Model:     c.Model,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Refs:      nil, // References would be loaded separately if needed
	}
}

// FromDomain converts domain model to Car
func FromDomain(car *model.Car) *Car {
	return &Car{
		ID:        car.ID,
		TenantID:  car.TenantID,
		Model:     car.Model,
		CreatedAt: car.CreatedAt,
		UpdatedAt: car.UpdatedAt,
	}
}

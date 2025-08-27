package dbmodel

import (
	"time"

	"github.com/jp-ryuji/go-sample/internal/domain/model"
	"gorm.io/gorm"
)

// Option represents the GORM model for Option
type Option struct {
	ID        string `gorm:"primaryKey;type:varchar(36);not null"`
	TenantID  string `gorm:"type:varchar(36);not null;index"`
	Name      string `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName specifies the table name for GORM
func (Option) TableName() string {
	return "options"
}

// ToDomain converts Option to domain model
func (o *Option) ToDomain() *model.Option {
	return &model.Option{
		ID:        o.ID,
		TenantID:  o.TenantID,
		Name:      o.Name,
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
		Refs:      nil, // References would be loaded separately if needed
	}
}

// FromDomain converts domain model to Option
func FromDomainOption(option *model.Option) *Option {
	return &Option{
		ID:        option.ID,
		TenantID:  option.TenantID,
		Name:      option.Name,
		CreatedAt: option.CreatedAt,
		UpdatedAt: option.UpdatedAt,
	}
}

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

	RentalOptions []RentalOption `gorm:"foreignKey:OptionID"`
}

// TableName specifies the table name for GORM
func (Option) TableName() string {
	return "options"
}

// OptionLoadOptions specifies which associations to load
type OptionLoadOptions struct {
	WithRentalOptions bool
}

// ToDomain converts Option to domain model with specified associations
func (o *Option) ToDomain(opts ...OptionLoadOptions) *model.Option {
	option := &model.Option{
		ID:        o.ID,
		TenantID:  o.TenantID,
		Name:      o.Name,
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
		Refs:      nil,
	}

	// If no options provided, return basic option
	if len(opts) == 0 {
		return option
	}

	opt := opts[0]

	// Only create Refs if rental options need to be loaded
	if opt.WithRentalOptions && len(o.RentalOptions) > 0 {
		rentalOptions := make(model.RentalOptions, len(o.RentalOptions))
		for i, rentalOption := range o.RentalOptions {
			rentalOptions[i] = rentalOption.ToDomain()
		}

		option.Refs = &model.OptionRefs{
			RentalOptions: rentalOptions,
		}
	}

	return option
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

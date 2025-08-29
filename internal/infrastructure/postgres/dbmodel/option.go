package dbmodel

import (
	"time"

	"github.com/jp-ryuji/go-sample/internal/domain/model"
)

// Option represents the database model for Option
type Option struct {
	ID        string     `json:"id"`
	TenantID  string     `json:"tenant_id"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	RentalOptions []RentalOption `json:"rental_options,omitempty"`
}

// TableName specifies the table name (kept for compatibility with existing code)
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

package dbmodel

import (
	"time"

	"github.com/jp-ryuji/go-arch-patterns/internal/domain/entity"
)

// Option represents the database model for Option
type Option struct {
	ID        string
	TenantID  string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	RentalOptions []RentalOption
}

// OptionLoadOptions specifies which associations to load
type OptionLoadOptions struct {
	WithRentalOptions bool
}

// ToDomain converts Option to domain model with specified associations
func (o *Option) ToDomain(opts ...OptionLoadOptions) *entity.Option {
	option := &entity.Option{
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
		rentalOptions := make(entity.RentalOptions, len(o.RentalOptions))
		for i, rentalOption := range o.RentalOptions {
			rentalOptions[i] = rentalOption.ToDomain()
		}

		option.Refs = &entity.OptionRefs{
			RentalOptions: rentalOptions,
		}
	}

	return option
}

// FromDomainOption converts domain model to Option
func FromDomainOption(option *entity.Option) *Option {
	return &Option{
		ID:        option.ID,
		TenantID:  option.TenantID,
		Name:      option.Name,
		CreatedAt: option.CreatedAt,
		UpdatedAt: option.UpdatedAt,
	}
}

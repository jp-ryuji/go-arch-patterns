package dbmodel

import (
	"time"

	"github.com/jp-ryuji/go-sample/internal/domain/model"
)

// Company represents the database model for Company
type Company struct {
	ID          string     `json:"id"`
	TenantID    string     `json:"tenant_id"`
	Name        string     `json:"name"`
	CompanySize string     `json:"company_size"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`

	Renters []Renter `json:"renters,omitempty"`
}

// CompanyLoadOptions specifies which associations to load
type CompanyLoadOptions struct {
	WithRenters bool
}

// ToDomain converts Company to domain model with specified associations
func (c *Company) ToDomain(opts ...CompanyLoadOptions) *model.Company {
	company := &model.Company{
		ID:          c.ID,
		TenantID:    c.TenantID,
		Name:        c.Name,
		CompanySize: model.NewCompanySize(c.CompanySize),
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
		Refs:        nil,
	}

	// If no options provided, return basic company
	if len(opts) == 0 {
		return company
	}

	option := opts[0]

	// Only create Refs if renters need to be loaded
	if option.WithRenters && len(c.Renters) > 0 {
		renters := make(model.Renters, len(c.Renters))
		for i, renter := range c.Renters {
			renters[i] = renter.ToDomain()
		}

		company.Refs = &model.CompanyRefs{
			Renters: renters,
		}
	}

	return company
}

// FromDomain converts domain model to Company
func FromDomainCompany(company *model.Company) *Company {
	return &Company{
		ID:          company.ID,
		TenantID:    company.TenantID,
		Name:        company.Name,
		CompanySize: company.CompanySize.String(),
		CreatedAt:   company.CreatedAt,
		UpdatedAt:   company.UpdatedAt,
	}
}

package dbmodel

import (
	"time"

	"github.com/jp-ryuji/go-sample/internal/domain/model"
)

// Company represents the database model for Company
type Company struct {
	ID          string
	RenterID    string
	TenantID    string
	Name        string
	CompanySize string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

// CompanyLoadOptions specifies which associations to load
type CompanyLoadOptions struct{}

// ToDomain converts Company to domain model with specified associations
func (c *Company) ToDomain(opts ...CompanyLoadOptions) *model.Company {
	company := &model.Company{
		ID:          c.ID,
		RenterID:    c.RenterID,
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

	// Create Refs if any associations need to be loaded
	company.Refs = &model.CompanyRefs{}

	return company
}

// FromDomain converts domain model to Company
func FromDomainCompany(company *model.Company) *Company {
	return &Company{
		ID:          company.ID,
		RenterID:    company.RenterID,
		TenantID:    company.TenantID,
		Name:        company.Name,
		CompanySize: company.CompanySize.String(),
		CreatedAt:   company.CreatedAt,
		UpdatedAt:   company.UpdatedAt,
	}
}

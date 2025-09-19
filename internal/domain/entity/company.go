package entity

import (
	"time"

	"github.com/oklog/ulid/v2"
)

// Companies is a slice of Company
type Companies []*Company

// Company represents a company entity (inherits from Renter)
type Company struct {
	ID          string
	RenterID    string
	TenantID    string
	Name        string
	CompanySize CompanySize
	CreatedAt   time.Time
	UpdatedAt   time.Time

	// References to related entities
	Refs *CompanyRefs
}

// CompanyRefs holds references to related entities
type CompanyRefs struct{}

// NewCompany creates a new Company
func NewCompany(renterID, tenantID, name string, companySize CompanySize, createdAt time.Time) *Company {
	return &Company{
		ID:          ulid.Make().String(),
		RenterID:    renterID,
		TenantID:    tenantID,
		Name:        name,
		CompanySize: companySize,
		CreatedAt:   createdAt,
		UpdatedAt:   createdAt,
	}
}

// WithRenterID creates a Company with a specific RenterID (for testing)
func (c *Company) WithRenterID(renterID string) *Company {
	c.RenterID = renterID
	return c
}

type CompanySize string

type CompanySizes []CompanySize

const (
	CompanySizeUnknown CompanySize = "unknown"
	CompanySizeSmall   CompanySize = "small_size_company"
	CompanySizeMedium  CompanySize = "medium_size_company"
	CompanySizeLarge   CompanySize = "large_size_company"
)

func NewCompanySize(s string) CompanySize {
	switch s {
	case CompanySizeSmall.String(),
		CompanySizeMedium.String(),
		CompanySizeLarge.String():
		return CompanySize(s)
	}
	return CompanySizeUnknown
}

func (m CompanySize) String() string {
	return string(m)
}

func (m CompanySize) Valid() bool {
	return m != CompanySizeUnknown && m != ""
}

func (ms CompanySizes) Slice() []string {
	dst := make([]string, 0, len(ms))
	for _, companySize := range ms {
		dst = append(dst, companySize.String())
	}
	return dst
}

func (ms CompanySizes) Valid() bool {
	for _, companySize := range ms {
		if !companySize.Valid() {
			return false
		}
	}
	return true
}

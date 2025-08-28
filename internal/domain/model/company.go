package model

import (
	"time"

	"github.com/google/uuid"
)

// Companies is a slice of Company
type Companies []*Company

// Company represents a company entity
type Company struct {
	ID          string
	TenantID    string
	Name        string
	CompanySize CompanySize
	CreatedAt   time.Time
	UpdatedAt   time.Time

	// References to related entities
	Refs *CompanyRefs
}

// CompanyRefs holds references to related entities
type CompanyRefs struct {
	Renters Renters
}

// NewCompany creates a new Company
func NewCompany(tenantID, name string, companySize CompanySize, createdAt time.Time) *Company {
	return &Company{
		ID:          uuid.New().String(),
		TenantID:    tenantID,
		Name:        name,
		CompanySize: companySize,
		CreatedAt:   createdAt,
		UpdatedAt:   createdAt,
	}
}

// WithID creates a Company with a specific ID (for testing)
func (c *Company) WithID(id string) *Company {
	c.ID = id
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

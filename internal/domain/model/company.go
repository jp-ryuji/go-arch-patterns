package model

import (
	"time"

	"github.com/jp-ryuji/go-sample/internal/pkg/id"
)

// CompanyRefs holds references to related entities for Company
type CompanyRefs struct {
	Tenant *Tenant
}

type Company struct {
	ID          string
	TenantID    string
	Name        string
	CompanySize CompanySize
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Refs *CompanyRefs
}

// Renter interface implementation
func (c *Company) GetID() string {
	return c.ID
}

func (c *Company) GetTenantID() string {
	return c.TenantID
}

func (c *Company) GetType() RenterModel {
	return RenterModelCompany
}

// End of Renter interface implementation

type Companys []*Company

func NewCompany(
	tenantID string,
	name string,
	companySize CompanySize,
	t time.Time,
) *Company {
	return &Company{
		ID:          id.New(),
		TenantID:    tenantID,
		Name:        name,
		CompanySize: companySize,
		CreatedAt:   t,
		UpdatedAt:   t,

		Refs: nil,
	}
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

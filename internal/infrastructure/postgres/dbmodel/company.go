package dbmodel

import (
	"time"

	"github.com/jp-ryuji/go-sample/internal/domain/model"
	"gorm.io/gorm"
)

// Company represents the GORM model for Company
type Company struct {
	ID          string `gorm:"primaryKey;type:varchar(36);not null"`
	TenantID    string `gorm:"type:varchar(36);not null;index"`
	Name        string `gorm:"type:varchar(255);not null"`
	CompanySize string `gorm:"type:varchar(50);not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// TableName specifies the table name for GORM
func (Company) TableName() string {
	return "companies"
}

// ToDomain converts Company to domain model
func (c *Company) ToDomain() *model.Company {
	return &model.Company{
		ID:          c.ID,
		TenantID:    c.TenantID,
		Name:        c.Name,
		CompanySize: model.NewCompanySize(c.CompanySize),
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
		Refs:        nil, // References would be loaded separately if needed
	}
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

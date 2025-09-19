package entity_test

import (
	"testing"

	"github.com/jp-ryuji/go-arch-patterns/internal/domain/entity"
	"github.com/jp-ryuji/go-arch-patterns/internal/domain/entity/factory"
)

// TestRenterCreation ensures that we can create renters
func TestRenterCreation(t *testing.T) {
	t.Parallel()

	// Create a company using the factory
	company := factory.NewCompany()

	// Create an individual using the factory
	individual, err := factory.NewIndividual()
	if err != nil {
		t.Fatalf("Failed to create individual: %v", err)
	}

	// Test company renter
	companyRenter := entity.NewRenter("tenant-123", entity.CompanyRenter, company.CreatedAt)
	if companyRenter.TenantID != "tenant-123" {
		t.Error("Company Renter TenantID mismatch")
	}
	if companyRenter.Type != entity.CompanyRenter {
		t.Error("Company Renter Type mismatch")
	}

	// Test individual renter
	individualRenter := entity.NewRenter("tenant-123", entity.IndividualRenter, individual.CreatedAt)
	if individualRenter.TenantID != "tenant-123" {
		t.Error("Individual Renter TenantID mismatch")
	}
	if individualRenter.Type != entity.IndividualRenter {
		t.Error("Individual Renter Type mismatch")
	}
}

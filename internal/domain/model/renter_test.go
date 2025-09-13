package model_test

import (
	"testing"

	"github.com/jp-ryuji/go-ddd/internal/domain/model"
	"github.com/jp-ryuji/go-ddd/internal/domain/model/factory"
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
	companyRenter := model.NewRenter("tenant-123", model.CompanyRenter, company.CreatedAt)
	if companyRenter.TenantID != "tenant-123" {
		t.Error("Company Renter TenantID mismatch")
	}
	if companyRenter.Type != model.CompanyRenter {
		t.Error("Company Renter Type mismatch")
	}

	// Test individual renter
	individualRenter := model.NewRenter("tenant-123", model.IndividualRenter, individual.CreatedAt)
	if individualRenter.TenantID != "tenant-123" {
		t.Error("Individual Renter TenantID mismatch")
	}
	if individualRenter.Type != model.IndividualRenter {
		t.Error("Individual Renter Type mismatch")
	}
}

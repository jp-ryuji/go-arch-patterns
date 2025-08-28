package model_test

import (
	"testing"

	"github.com/jp-ryuji/go-sample/internal/domain/model"
	"github.com/jp-ryuji/go-sample/internal/domain/model/factory"
)

// TestRenterCreation ensures that we can create renters with different entity types
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
	companyRenter := model.NewRenter("tenant-123", company.ID, model.CompanyRenter, company.CreatedAt)
	if companyRenter.RenterEntityID != company.ID {
		t.Error("Company RenterEntityID mismatch")
	}
	if companyRenter.RenterEntityType != model.CompanyRenter {
		t.Error("Company RenterEntityType mismatch")
	}

	// Test individual renter
	individualRenter := model.NewRenter("tenant-123", individual.ID, model.IndividualRenter, individual.CreatedAt)
	if individualRenter.RenterEntityID != individual.ID {
		t.Error("Individual RenterEntityID mismatch")
	}
	if individualRenter.RenterEntityType != model.IndividualRenter {
		t.Error("Individual RenterEntityType mismatch")
	}
}

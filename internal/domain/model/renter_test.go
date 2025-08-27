package model_test

import (
	"testing"

	"github.com/jp-ryuji/go-sample/internal/domain/model"
	"github.com/jp-ryuji/go-sample/internal/domain/model/factory"
)

// TestRenterInterface ensures that both Company and Individual implement the RenterEntityInstance interface
func TestRenterInterface(t *testing.T) {
	t.Parallel()

	// Create a company using the factory
	company := factory.NewCompany()

	// Create an individual using the factory
	individual, err := factory.NewIndividual()
	if err != nil {
		t.Fatalf("Failed to create individual: %v", err)
	}

	// Verify that both implement the RenterEntityInstance interface
	var companyRenter model.RenterEntityInstance = company
	var individualRenter model.RenterEntityInstance = individual

	// Test company
	if companyRenter.GetID() != company.ID {
		t.Error("Company GetID() mismatch")
	}
	if companyRenter.GetEntityType() != model.RenterEntityCompany {
		t.Error("Company GetEntityType() mismatch")
	}

	// Test individual
	if individualRenter.GetID() != individual.ID {
		t.Error("Individual GetID() mismatch")
	}
	if individualRenter.GetEntityType() != model.RenterEntityIndividual {
		t.Error("Individual GetEntityType() mismatch")
	}
}

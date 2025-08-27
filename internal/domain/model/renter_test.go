package model_test

import (
	"testing"

	"github.com/jp-ryuji/go-sample/internal/domain/model"
	"github.com/jp-ryuji/go-sample/internal/domain/model/factory"
)

// TestRenterInterface ensures that both Company and Individual implement the ConcreteRenter interface
func TestRenterInterface(t *testing.T) {
	t.Parallel()

	// Create a company using the factory
	company := factory.NewCompany()

	// Create an individual using the factory
	individual, err := factory.NewIndividual()
	if err != nil {
		t.Fatalf("Failed to create individual: %v", err)
	}

	// Verify that both implement the ConcreteRenter interface
	var companyRenter model.ConcreteRenter = company
	var individualRenter model.ConcreteRenter = individual

	// Test company
	if companyRenter.GetID() != company.ID {
		t.Error("Company GetID() mismatch")
	}
	if companyRenter.GetTenantID() != company.TenantID {
		t.Error("Company GetTenantID() mismatch")
	}
	if companyRenter.GetType() != model.RenterModelCompany {
		t.Error("Company GetType() mismatch")
	}

	// Test individual
	if individualRenter.GetID() != individual.ID {
		t.Error("Individual GetID() mismatch")
	}
	if individualRenter.GetTenantID() != individual.TenantID {
		t.Error("Individual GetTenantID() mismatch")
	}
	if individualRenter.GetType() != model.RenterModelIndividual {
		t.Error("Individual GetType() mismatch")
	}
}

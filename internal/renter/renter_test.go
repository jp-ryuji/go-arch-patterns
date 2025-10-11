package renter_test

import (
	"testing"
	"time"

	renter "github.com/jp-ryuji/go-arch-patterns/internal/renter"
)

func TestRenterCreation(t *testing.T) {
	t.Parallel()

	companyRenter := renter.NewRenter("tenant-123", renter.TypeCompany, time.Now())
	if companyRenter.TenantID != "tenant-123" {
		t.Error("Company renter tenant mismatch")
	}
	if companyRenter.Type != renter.TypeCompany {
		t.Error("Company renter type mismatch")
	}

	individualRenter := renter.NewRenter("tenant-123", renter.TypeIndividual, time.Now())
	if individualRenter.TenantID != "tenant-123" {
		t.Error("Individual renter tenant mismatch")
	}
	if individualRenter.Type != renter.TypeIndividual {
		t.Error("Individual renter type mismatch")
	}
}

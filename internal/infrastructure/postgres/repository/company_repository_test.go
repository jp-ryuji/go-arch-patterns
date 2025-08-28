//go:build integration

// Package repository_test provides integration tests for the repository implementations.
package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/jp-ryuji/go-sample/internal/domain/model"
	"github.com/jp-ryuji/go-sample/internal/domain/repository"
	companyrepo "github.com/jp-ryuji/go-sample/internal/infrastructure/postgres/repository"
	"github.com/jp-ryuji/go-sample/internal/infrastructure/postgres/repository/testutil"
	"github.com/stretchr/testify/require"
)

// testCompanySetup is a helper function that provides common test setup for company tests
func testCompanySetup(t *testing.T, tenantCode string) (repository.CompanyRepository, context.Context, *model.Tenant) {
	t.Helper()

	// Skip this test if not running integration tests
	testutil.SkipIfShort(t)

	repo := companyrepo.NewCompanyRepository()
	ctx := context.Background()
	tenant := testutil.CreateTestTenant(t, tenantCode)

	return repo, ctx, tenant
}

// TestCompanyRepository_Create tests the Create method of the company repository.
func TestCompanyRepository_Create(t *testing.T) {
	repo, ctx, tenant := testCompanySetup(t, "test-tenant-company-create")

	// Create a company
	company := model.NewCompany(tenant.ID, "Test Company", model.CompanySizeMedium, time.Now())
	err := repo.Create(ctx, company)
	require.NoError(t, err)

	// Verify the company was created
	foundCompany, err := repo.GetByID(ctx, company.ID)
	require.NoError(t, err)
	require.Equal(t, company.ID, foundCompany.ID)
	require.Equal(t, company.TenantID, foundCompany.TenantID)
	require.Equal(t, company.Name, foundCompany.Name)
	require.Equal(t, company.CompanySize, foundCompany.CompanySize)
}

// TestCompanyRepository_GetByID tests the GetByID method of the company repository.
func TestCompanyRepository_GetByID(t *testing.T) {
	repo, ctx, tenant := testCompanySetup(t, "test-tenant-company-get")

	// Create a company
	company := model.NewCompany(tenant.ID, "Another Test Company", model.CompanySizeLarge, time.Now())
	err := repo.Create(ctx, company)
	require.NoError(t, err)

	// Get the company by ID
	foundCompany, err := repo.GetByID(ctx, company.ID)
	require.NoError(t, err)
	require.Equal(t, company.ID, foundCompany.ID)
	require.Equal(t, company.TenantID, foundCompany.TenantID)
	require.Equal(t, company.Name, foundCompany.Name)
	require.Equal(t, company.CompanySize, foundCompany.CompanySize)
}

// TestCompanyRepository_GetByID_NotFound tests the GetByID method when a company doesn't exist.
func TestCompanyRepository_GetByID_NotFound(t *testing.T) {
	repo, ctx, _ := testCompanySetup(t, "test-tenant-company-get-not-found")

	// Try to get a company that doesn't exist
	_, err := repo.GetByID(ctx, "non-existent-id")
	require.Error(t, err)
}

// TestCompanyRepository_Update tests the Update method of the company repository.
func TestCompanyRepository_Update(t *testing.T) {
	repo, ctx, tenant := testCompanySetup(t, "test-tenant-company-update")

	// Create a company
	company := model.NewCompany(tenant.ID, "Original Company", model.CompanySizeSmall, time.Now())
	err := repo.Create(ctx, company)
	require.NoError(t, err)

	// Update the company
	originalUpdatedAt := company.UpdatedAt
	company.Name = "Updated Company"
	company.CompanySize = model.CompanySizeLarge
	company.UpdatedAt = time.Now()
	err = repo.Update(ctx, company)
	require.NoError(t, err)

	// Verify the update
	updatedCompany, err := repo.GetByID(ctx, company.ID)
	require.NoError(t, err)
	require.Equal(t, "Updated Company", updatedCompany.Name)
	require.Equal(t, model.CompanySizeLarge, updatedCompany.CompanySize)
	require.True(t, updatedCompany.UpdatedAt.After(originalUpdatedAt))
}

// TestCompanyRepository_Delete tests the Delete method of the company repository.
func TestCompanyRepository_Delete(t *testing.T) {
	repo, ctx, tenant := testCompanySetup(t, "test-tenant-company-delete")

	// Create a company
	company := model.NewCompany(tenant.ID, "Company to Delete", model.CompanySizeSmall, time.Now())
	err := repo.Create(ctx, company)
	require.NoError(t, err)

	// Delete the company
	err = repo.Delete(ctx, company.ID)
	require.NoError(t, err)

	// Verify the company is deleted
	_, err = repo.GetByID(ctx, company.ID)
	require.Error(t, err)
}

// TestCompanyRepository_Delete_NotFound tests the Delete method when a company doesn't exist.
func TestCompanyRepository_Delete_NotFound(t *testing.T) {
	repo, ctx, _ := testCompanySetup(t, "test-tenant-company-delete-not-found")

	// Try to delete a company that doesn't exist
	err := repo.Delete(ctx, "non-existent-id")
	require.NoError(t, err) // Delete should be idempotent
}

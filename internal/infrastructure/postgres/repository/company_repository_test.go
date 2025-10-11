//go:build integration

// Package repository_test provides integration tests for the repository implementations.
package repository_test

import (
	"context"
	"testing"
	"time"

	companyrepo "github.com/jp-ryuji/go-arch-patterns/internal/infrastructure/postgres/repository"
	"github.com/jp-ryuji/go-arch-patterns/internal/infrastructure/postgres/repository/testutil"
	renter "github.com/jp-ryuji/go-arch-patterns/internal/renter"
	tenant "github.com/jp-ryuji/go-arch-patterns/internal/tenant"
	"github.com/stretchr/testify/require"
)

func testCompanySetup(t *testing.T, tenantCode string) (renter.CompanyRepository, renter.RenterRepository, context.Context, *tenant.Tenant) {
	t.Helper()
	testutil.SkipIfShort(t)

	repo := companyrepo.NewCompanyRepository(testutil.DBClient)
	renterRepo := companyrepo.NewRenterRepository(testutil.DBClient)
	ctx := context.Background()
	tenantModel := testutil.CreateTestTenant(t, tenantCode)

	return repo, renterRepo, ctx, tenantModel
}

func TestCompanyRepository_Create(t *testing.T) {
	repo, renterRepo, ctx, tenantModel := testCompanySetup(t, "test-tenant-company-create")

	renterModel := renter.NewRenter(tenantModel.ID, renter.TypeCompany, time.Now())
	err := renterRepo.Create(ctx, renterModel)
	require.NoError(t, err)

	companyModel := renter.NewCompany(renterModel.ID, tenantModel.ID, "Test Company", renter.CompanySizeMedium, time.Now())
	err = repo.Create(ctx, companyModel)
	require.NoError(t, err)

	foundCompany, err := repo.GetByID(ctx, companyModel.RenterID)
	require.NoError(t, err)
	require.Equal(t, companyModel.RenterID, foundCompany.RenterID)
	require.Equal(t, companyModel.TenantID, foundCompany.TenantID)
	require.Equal(t, companyModel.Name, foundCompany.Name)
	require.Equal(t, companyModel.CompanySize, foundCompany.CompanySize)
}

func TestCompanyRepository_GetByID(t *testing.T) {
	repo, renterRepo, ctx, tenantModel := testCompanySetup(t, "test-tenant-company-get")

	renterModel := renter.NewRenter(tenantModel.ID, renter.TypeCompany, time.Now())
	err := renterRepo.Create(ctx, renterModel)
	require.NoError(t, err)

	companyModel := renter.NewCompany(renterModel.ID, tenantModel.ID, "Another Test Company", renter.CompanySizeLarge, time.Now())
	err = repo.Create(ctx, companyModel)
	require.NoError(t, err)

	foundCompany, err := repo.GetByID(ctx, companyModel.RenterID)
	require.NoError(t, err)
	require.Equal(t, companyModel.RenterID, foundCompany.RenterID)
	require.Equal(t, companyModel.TenantID, foundCompany.TenantID)
	require.Equal(t, companyModel.Name, foundCompany.Name)
	require.Equal(t, companyModel.CompanySize, foundCompany.CompanySize)
}

func TestCompanyRepository_GetByID_NotFound(t *testing.T) {
	repo, _, ctx, _ := testCompanySetup(t, "test-tenant-company-get-not-found")
	_, err := repo.GetByID(ctx, "non-existent-id")
	require.Error(t, err)
}

func TestCompanyRepository_Update(t *testing.T) {
	repo, renterRepo, ctx, tenantModel := testCompanySetup(t, "test-tenant-company-update")

	renterModel := renter.NewRenter(tenantModel.ID, renter.TypeCompany, time.Now())
	err := renterRepo.Create(ctx, renterModel)
	require.NoError(t, err)

	companyModel := renter.NewCompany(renterModel.ID, tenantModel.ID, "Original Company", renter.CompanySizeSmall, time.Now())
	err = repo.Create(ctx, companyModel)
	require.NoError(t, err)

	originalUpdatedAt := companyModel.UpdatedAt
	companyModel.Name = "Updated Company"
	companyModel.CompanySize = renter.CompanySizeLarge
	companyModel.UpdatedAt = time.Now()
	err = repo.Update(ctx, companyModel)
	require.NoError(t, err)

	updatedCompany, err := repo.GetByID(ctx, companyModel.RenterID)
	require.NoError(t, err)
	require.Equal(t, "Updated Company", updatedCompany.Name)
	require.Equal(t, renter.CompanySizeLarge, updatedCompany.CompanySize)
	require.True(t, updatedCompany.UpdatedAt.After(originalUpdatedAt))
}

func TestCompanyRepository_Delete(t *testing.T) {
	repo, renterRepo, ctx, tenantModel := testCompanySetup(t, "test-tenant-company-delete")

	renterModel := renter.NewRenter(tenantModel.ID, renter.TypeCompany, time.Now())
	err := renterRepo.Create(ctx, renterModel)
	require.NoError(t, err)

	companyModel := renter.NewCompany(renterModel.ID, tenantModel.ID, "Company to Delete", renter.CompanySizeSmall, time.Now())
	err = repo.Create(ctx, companyModel)
	require.NoError(t, err)

	err = repo.Delete(ctx, companyModel.RenterID)
	require.NoError(t, err)

	_, err = repo.GetByID(ctx, companyModel.RenterID)
	require.Error(t, err)
}

func TestCompanyRepository_Delete_NotFound(t *testing.T) {
	repo, _, ctx, _ := testCompanySetup(t, "test-tenant-company-delete-not-found")
	err := repo.Delete(ctx, "non-existent-id")
	require.NoError(t, err)
}

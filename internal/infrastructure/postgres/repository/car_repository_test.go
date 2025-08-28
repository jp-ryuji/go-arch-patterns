//go:build integration

// Package repository_test provides integration tests for the repository implementations.
package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/jp-ryuji/go-sample/internal/domain/model"
	"github.com/jp-ryuji/go-sample/internal/domain/repository"
	carrepo "github.com/jp-ryuji/go-sample/internal/infrastructure/postgres/repository"
	"github.com/jp-ryuji/go-sample/internal/infrastructure/postgres/repository/testutil"
	"github.com/stretchr/testify/require"
)

// testSetup is a helper function that provides common test setup
func testSetup(t *testing.T, tenantCode string) (repository.CarRepository, context.Context, *model.Tenant) {
	t.Helper()

	// Skip this test if not running integration tests
	testutil.SkipIfShort(t)

	repo := carrepo.NewCarRepository()
	ctx := context.Background()
	tenant := testutil.CreateTestTenant(t, tenantCode)

	return repo, ctx, tenant
}

// TestCarRepository_Create tests the Create method of the car repository.
func TestCarRepository_Create(t *testing.T) {
	repo, ctx, tenant := testSetup(t, "test-tenant-create")

	// Create a car
	car := model.NewCar(tenant.ID, "HARRIER", time.Now())
	err := repo.Create(ctx, car)
	require.NoError(t, err)

	// Verify the car was created
	foundCar, err := repo.GetByID(ctx, car.ID)
	require.NoError(t, err)
	require.Equal(t, car.ID, foundCar.ID)
	require.Equal(t, car.TenantID, foundCar.TenantID)
	require.Equal(t, car.Model, foundCar.Model)
}

// TestCarRepository_GetByID tests the GetByID method of the car repository.
func TestCarRepository_GetByID(t *testing.T) {
	repo, ctx, tenant := testSetup(t, "test-tenant-get")

	// Create a car
	car := model.NewCar(tenant.ID, "RAV4", time.Now())
	err := repo.Create(ctx, car)
	require.NoError(t, err)

	// Get the car by ID
	foundCar, err := repo.GetByID(ctx, car.ID)
	require.NoError(t, err)
	require.Equal(t, car.ID, foundCar.ID)
	require.Equal(t, car.TenantID, foundCar.TenantID)
	require.Equal(t, car.Model, foundCar.Model)
}

// TestCarRepository_GetByID_NotFound tests the GetByID method when a car doesn't exist.
func TestCarRepository_GetByID_NotFound(t *testing.T) {
	repo, ctx, _ := testSetup(t, "test-tenant-get-not-found")

	// Try to get a car that doesn't exist
	_, err := repo.GetByID(ctx, "non-existent-id")
	require.Error(t, err)
}

// TestCarRepository_Update tests the Update method of the car repository.
func TestCarRepository_Update(t *testing.T) {
	repo, ctx, tenant := testSetup(t, "test-tenant-update")

	// Create a car
	car := model.NewCar(tenant.ID, "HR-V", time.Now())
	err := repo.Create(ctx, car)
	require.NoError(t, err)

	// Update the car
	originalUpdatedAt := car.UpdatedAt
	car.Model = "CR-V"
	car.UpdatedAt = time.Now()
	err = repo.Update(ctx, car)
	require.NoError(t, err)

	// Verify the update
	updatedCar, err := repo.GetByID(ctx, car.ID)
	require.NoError(t, err)
	require.Equal(t, "CR-V", updatedCar.Model)
	require.True(t, updatedCar.UpdatedAt.After(originalUpdatedAt))
}

// TestCarRepository_Delete tests the Delete method of the car repository.
func TestCarRepository_Delete(t *testing.T) {
	repo, ctx, tenant := testSetup(t, "test-tenant-delete")

	// Create a car
	car := model.NewCar(tenant.ID, "Tesla Model Y", time.Now())
	err := repo.Create(ctx, car)
	require.NoError(t, err)

	// Delete the car
	err = repo.Delete(ctx, car.ID)
	require.NoError(t, err)

	// Verify the car is deleted
	_, err = repo.GetByID(ctx, car.ID)
	require.Error(t, err)
}

// TestCarRepository_Delete_NotFound tests the Delete method when a car doesn't exist.
func TestCarRepository_Delete_NotFound(t *testing.T) {
	repo, ctx, _ := testSetup(t, "test-tenant-delete-not-found")

	// Try to delete a car that doesn't exist
	err := repo.Delete(ctx, "non-existent-id")
	require.NoError(t, err) // Delete should be idempotent
}

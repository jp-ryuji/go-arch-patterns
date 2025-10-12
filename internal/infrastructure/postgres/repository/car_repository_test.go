//go:build integration

package repository_test

import (
	"context"
	"testing"
	"time"

	carrepo "github.com/jp-ryuji/go-arch-patterns/internal/infrastructure/postgres/repository"
	"github.com/jp-ryuji/go-arch-patterns/internal/infrastructure/postgres/repository/testutil"
	rental "github.com/jp-ryuji/go-arch-patterns/internal/rental"
	tenant "github.com/jp-ryuji/go-arch-patterns/internal/tenant"
	"github.com/stretchr/testify/require"
)

func testSetup(t *testing.T, tenantCode string) (rental.CarRepository, context.Context, *tenant.Tenant) {
	t.Helper()
	testutil.SkipIfShort(t)

	repo := carrepo.NewCarRepository(testutil.DBClient)
	ctx := context.Background()
	tenantModel := testutil.CreateTestTenant(t, tenantCode)

	return repo, ctx, tenantModel
}

func TestCarRepository_Create(t *testing.T) {
	repo, ctx, tenantModel := testSetup(t, "test-tenant-create")

	car := rental.NewCar(tenantModel.ID, "HARRIER", time.Now())
	err := repo.Create(ctx, car)
	require.NoError(t, err)

	foundCar, err := repo.GetByID(ctx, car.ID)
	require.NoError(t, err)
	require.Equal(t, car.ID, foundCar.ID)
	require.Equal(t, car.TenantID, foundCar.TenantID)
	require.Equal(t, car.Model, foundCar.Model)
}

func TestCarRepository_GetByID(t *testing.T) {
	repo, ctx, tenantModel := testSetup(t, "test-tenant-get")

	car := rental.NewCar(tenantModel.ID, "RAV4", time.Now())
	err := repo.Create(ctx, car)
	require.NoError(t, err)

	foundCar, err := repo.GetByID(ctx, car.ID)
	require.NoError(t, err)
	require.Equal(t, car.ID, foundCar.ID)
	require.Equal(t, car.TenantID, foundCar.TenantID)
	require.Equal(t, car.Model, foundCar.Model)
}

func TestCarRepository_GetByIDWithTenant(t *testing.T) {
	repo, ctx, tenantModel := testSetup(t, "test-tenant-get-with-tenant")

	car := rental.NewCar(tenantModel.ID, "Prius", time.Now())
	err := repo.Create(ctx, car)
	require.NoError(t, err)

	foundCar, err := repo.GetByIDWithTenant(ctx, car.ID)
	require.NoError(t, err)
	require.Equal(t, car.ID, foundCar.ID)
	require.Equal(t, car.TenantID, foundCar.TenantID)
	require.Equal(t, car.Model, foundCar.Model)
}

func TestCarRepository_GetByID_NotFound(t *testing.T) {
	repo, ctx, _ := testSetup(t, "test-tenant-get-not-found")
	_, err := repo.GetByID(ctx, "non-existent-id")
	require.Error(t, err)
}

func TestCarRepository_Update(t *testing.T) {
	repo, ctx, tenantModel := testSetup(t, "test-tenant-update")

	car := rental.NewCar(tenantModel.ID, "HR-V", time.Now())
	err := repo.Create(ctx, car)
	require.NoError(t, err)

	originalUpdatedAt := car.UpdatedAt
	car.Model = "CR-V"
	car.UpdatedAt = time.Now()
	err = repo.Update(ctx, car)
	require.NoError(t, err)

	updatedCar, err := repo.GetByID(ctx, car.ID)
	require.NoError(t, err)
	require.Equal(t, "CR-V", updatedCar.Model)
	require.True(t, updatedCar.UpdatedAt.After(originalUpdatedAt))
}

func TestCarRepository_Delete(t *testing.T) {
	repo, ctx, tenantModel := testSetup(t, "test-tenant-delete")

	car := rental.NewCar(tenantModel.ID, "Tesla Model Y", time.Now())
	err := repo.Create(ctx, car)
	require.NoError(t, err)

	err = repo.Delete(ctx, car.ID)
	require.NoError(t, err)

	_, err = repo.GetByID(ctx, car.ID)
	require.Error(t, err)
}

func TestCarRepository_Delete_NotFound(t *testing.T) {
	repo, ctx, _ := testSetup(t, "test-tenant-delete-not-found")
	err := repo.Delete(ctx, "non-existent-id")
	require.NoError(t, err)
}

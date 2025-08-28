//go:build integration

package dbmodel_test

import (
	"os"
	"testing"
	"time"

	"github.com/jp-ryuji/go-sample/internal/domain/model"
	"github.com/jp-ryuji/go-sample/internal/domain/model/factory"
	"github.com/jp-ryuji/go-sample/internal/infrastructure/postgres/dbmodel"
	"github.com/jp-ryuji/go-sample/internal/infrastructure/postgres/repository/testutil"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	// Initialize the test environment
	if err := testutil.SetupTestEnvironment(); err != nil {
		panic(err)
	}

	// Run the tests
	code := m.Run()

	// Clean up the test environment
	if err := testutil.TeardownTestEnvironment(); err != nil {
		panic(err)
	}

	os.Exit(code)
}

func TestCarToDomainWithTenant(t *testing.T) {
	testutil.SkipIfShort(t)

	db := testutil.DBClient.DB

	// Create a tenant
	tenant := factory.NewTenantWithCode("test-tenant")
	tenantDB := dbmodel.FromDomainTenant(tenant)
	err := db.Create(tenantDB).Error
	require.NoError(t, err)

	// Create a car associated with the tenant
	car := model.NewCar(tenant.ID, "Toyota Camry", time.Now())
	carDB := dbmodel.FromDomainCar(car)
	err = db.Create(carDB).Error
	require.NoError(t, err)

	// Test 1: Basic ToDomain conversion (without eager loading)
	t.Run("ToDomain without eager loading", func(t *testing.T) {
		var fetchedCarDB dbmodel.Car
		err = db.Where("id = ?", car.ID).First(&fetchedCarDB).Error
		require.NoError(t, err)

		domainCar := fetchedCarDB.ToDomain()
		require.Equal(t, car.ID, domainCar.ID)
		require.Equal(t, car.TenantID, domainCar.TenantID)
		require.Equal(t, car.Model, domainCar.Model)
		// Tenant reference should be nil since we didn't eager load it
		require.Nil(t, domainCar.Refs)
	})

	// Test 2: ToDomainWithTenant conversion (with eager loading)
	t.Run("ToDomainWithTenant with eager loading", func(t *testing.T) {
		var fetchedCarDB dbmodel.Car
		err = db.Preload("Tenant").Where("id = ?", car.ID).First(&fetchedCarDB).Error
		require.NoError(t, err)

		domainCar := fetchedCarDB.ToDomain(dbmodel.CarLoadOptions{WithTenant: true})
		require.Equal(t, car.ID, domainCar.ID)
		require.Equal(t, car.TenantID, domainCar.TenantID)
		require.Equal(t, car.Model, domainCar.Model)
		// Tenant reference should be loaded
		require.NotNil(t, domainCar.Refs)
		require.NotNil(t, domainCar.Refs.Tenant)
		require.Equal(t, tenant.ID, domainCar.Refs.Tenant.ID)
		require.Equal(t, tenant.Code, domainCar.Refs.Tenant.Code)
	})

	// Test 3: ToDomainWithTenant when tenant is not found
	t.Run("ToDomainWithTenant with missing tenant", func(t *testing.T) {
		// Create a car associated with a tenant
		car := model.NewCar(tenant.ID, "Honda Civic", time.Now())
		carDB := dbmodel.FromDomainCar(car)
		err = db.Create(carDB).Error
		require.NoError(t, err)

		// Fetch the car without preloading the tenant
		var fetchedCarDB dbmodel.Car
		err = db.Where("id = ?", car.ID).First(&fetchedCarDB).Error
		require.NoError(t, err)

		domainCar := fetchedCarDB.ToDomain(dbmodel.CarLoadOptions{WithTenant: true})
		require.Equal(t, car.ID, domainCar.ID)
		require.Equal(t, car.TenantID, domainCar.TenantID)
		require.Equal(t, car.Model, domainCar.Model)
		// Tenant reference should not be loaded since we didn't preload it
		require.Nil(t, domainCar.Refs)
	})
}

func TestCarRepositoryWithTenant(t *testing.T) {
	testutil.SkipIfShort(t)

	db := testutil.DBClient.DB

	// Create a tenant
	tenant := factory.NewTenantWithCode("car-repo-tenant")
	tenantDB := dbmodel.FromDomainTenant(tenant)
	err := db.Create(tenantDB).Error
	require.NoError(t, err)

	// Create a car associated with the tenant
	car := model.NewCar(tenant.ID, "Nissan Altima", time.Now())
	carDB := dbmodel.FromDomainCar(car)
	err = db.Create(carDB).Error
	require.NoError(t, err)

	// Test fetching car with tenant using joins
	t.Run("Fetch car with tenant using joins", func(t *testing.T) {
		var fetchedCarDB dbmodel.Car
		err = db.Joins("Tenant").Where("cars.id = ?", car.ID).First(&fetchedCarDB).Error
		require.NoError(t, err)

		domainCar := fetchedCarDB.ToDomain(dbmodel.CarLoadOptions{WithTenant: true})
		require.Equal(t, car.ID, domainCar.ID)
		require.Equal(t, car.TenantID, domainCar.TenantID)
		require.Equal(t, car.Model, domainCar.Model)
		// Tenant reference should be loaded
		require.NotNil(t, domainCar.Refs)
		require.NotNil(t, domainCar.Refs.Tenant)
		require.Equal(t, tenant.ID, domainCar.Refs.Tenant.ID)
		require.Equal(t, tenant.Code, domainCar.Refs.Tenant.Code)
	})
}

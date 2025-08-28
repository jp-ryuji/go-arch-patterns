//go:build integration

package dbmodel_test

import (
	"testing"
	"time"

	"github.com/jp-ryuji/go-sample/internal/domain/model"
	"github.com/jp-ryuji/go-sample/internal/domain/model/factory"
	"github.com/jp-ryuji/go-sample/internal/infrastructure/postgres/dbmodel"
	"github.com/jp-ryuji/go-sample/internal/infrastructure/postgres/repository/testutil"
	"github.com/stretchr/testify/require"
)

func TestRentalToDomain(t *testing.T) {
	testutil.SkipIfShort(t)

	db := testutil.DBClient.DB

	// Create a tenant
	tenant := factory.NewTenantWithCode("rental-test-tenant")
	tenantDB := dbmodel.FromDomainTenant(tenant)
	err := db.Create(tenantDB).Error
	require.NoError(t, err)

	// Create a company
	company := factory.NewCompany()
	companyDB := dbmodel.FromDomainCompany(company)
	err = db.Create(companyDB).Error
	require.NoError(t, err)

	// Create a renter that references the company
	renter := model.NewRenter(tenant.ID, company.ID, model.CompanyRenter, time.Now())
	renterDB := dbmodel.FromDomainRenter(renter)
	err = db.Create(renterDB).Error
	require.NoError(t, err)

	// Create a car
	car := model.NewCar(tenant.ID, "Toyota Camry", time.Now())
	carDB := dbmodel.FromDomainCar(car)
	err = db.Create(carDB).Error
	require.NoError(t, err)

	// Create a rental
	now := time.Now()
	rental := model.NewRental(tenant.ID, car.ID, renter.ID, now.AddDate(0, 0, 1), now.AddDate(0, 0, 7))
	rentalDB := dbmodel.FromDomainRental(rental)
	err = db.Create(rentalDB).Error
	require.NoError(t, err)

	// Test 1: Basic ToDomain conversion (without eager loading)
	t.Run("ToDomain without eager loading", func(t *testing.T) {
		var fetchedRentalDB dbmodel.Rental
		err = db.Where("id = ?", rental.ID).First(&fetchedRentalDB).Error
		require.NoError(t, err)

		domainRental := fetchedRentalDB.ToDomain()
		require.Equal(t, rental.ID, domainRental.ID)
		require.Equal(t, rental.TenantID, domainRental.TenantID)
		require.Equal(t, rental.CarID, domainRental.CarID)
		require.Equal(t, rental.RenterID, domainRental.RenterID)
		require.Equal(t, rental.StartsAt.Unix(), domainRental.StartsAt.Unix())
		require.Equal(t, rental.EndsAt.Unix(), domainRental.EndsAt.Unix())
		// References should be nil since we didn't eager load them
		require.Nil(t, domainRental.Refs)
	})

	// Test 2: ToDomain with all associations
	t.Run("ToDomain with all associations", func(t *testing.T) {
		var fetchedRentalDB dbmodel.Rental
		err = db.Preload("Tenant").Preload("Car").Preload("Renter").Where("id = ?", rental.ID).First(&fetchedRentalDB).Error
		require.NoError(t, err)

		domainRental := fetchedRentalDB.ToDomain(dbmodel.RentalLoadOptions{
			WithTenant: true,
			WithCar:    true,
			WithRenter: true,
		})
		require.Equal(t, rental.ID, domainRental.ID)
		require.Equal(t, rental.TenantID, domainRental.TenantID)
		require.Equal(t, rental.CarID, domainRental.CarID)
		require.Equal(t, rental.RenterID, domainRental.RenterID)
		require.Equal(t, rental.StartsAt.Unix(), domainRental.StartsAt.Unix())
		require.Equal(t, rental.EndsAt.Unix(), domainRental.EndsAt.Unix())
		// All references should be loaded
		require.NotNil(t, domainRental.Refs)
		require.NotNil(t, domainRental.Refs.Tenant)
		require.NotNil(t, domainRental.Refs.Car)
		require.NotNil(t, domainRental.Refs.Renter)
		require.Equal(t, tenant.ID, domainRental.Refs.Tenant.ID)
		require.Equal(t, car.ID, domainRental.Refs.Car.ID)
		require.Equal(t, renter.ID, domainRental.Refs.Renter.ID)
	})

	// Test 3: ToDomain with partial associations
	t.Run("ToDomain with partial associations", func(t *testing.T) {
		var fetchedRentalDB dbmodel.Rental
		err = db.Preload("Tenant").Preload("Car").Where("id = ?", rental.ID).First(&fetchedRentalDB).Error
		require.NoError(t, err)

		domainRental := fetchedRentalDB.ToDomain(dbmodel.RentalLoadOptions{
			WithTenant: true,
			WithCar:    true,
			WithRenter: false, // Don't load renter
		})
		require.Equal(t, rental.ID, domainRental.ID)
		require.Equal(t, rental.TenantID, domainRental.TenantID)
		require.Equal(t, rental.CarID, domainRental.CarID)
		require.Equal(t, rental.RenterID, domainRental.RenterID)
		require.Equal(t, rental.StartsAt.Unix(), domainRental.StartsAt.Unix())
		require.Equal(t, rental.EndsAt.Unix(), domainRental.EndsAt.Unix())
		// Only tenant and car should be loaded
		require.NotNil(t, domainRental.Refs)
		require.NotNil(t, domainRental.Refs.Tenant)
		require.NotNil(t, domainRental.Refs.Car)
		require.Nil(t, domainRental.Refs.Renter) // Renter should be nil
		require.Equal(t, tenant.ID, domainRental.Refs.Tenant.ID)
		require.Equal(t, car.ID, domainRental.Refs.Car.ID)
	})
}

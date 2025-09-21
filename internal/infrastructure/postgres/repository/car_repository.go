package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jp-ryuji/go-arch-patterns/internal/domain/entity"
	"github.com/jp-ryuji/go-arch-patterns/internal/domain/repository"
	"github.com/jp-ryuji/go-arch-patterns/internal/infrastructure/postgres/entgen"
	car "github.com/jp-ryuji/go-arch-patterns/internal/infrastructure/postgres/entgen/car"
)

type carRepository struct {
	client *entgen.Client
}

// NewCarRepository creates a new car repository
func NewCarRepository(client *entgen.Client) repository.CarRepository {
	return &carRepository{
		client: client,
	}
}

// Create inserts a new car into the database
func (r *carRepository) Create(ctx context.Context, car *entity.Car) error {
	_, err := r.client.Car.
		Create().
		SetID(car.ID).
		SetTenantID(car.TenantID).
		SetModel(car.Model).
		Save(ctx)
	return err
}

// CreateInTx inserts a new car into the database within a transaction
func (r *carRepository) CreateInTx(ctx context.Context, tx *entgen.Tx, car *entity.Car) error {
	_, err := tx.Car.
		Create().
		SetID(car.ID).
		SetTenantID(car.TenantID).
		SetModel(car.Model).
		Save(ctx)
	return err
}

// GetByID retrieves a car by its ID
func (r *carRepository) GetByID(ctx context.Context, id string) (*entity.Car, error) {
	carDB, err := r.client.Car.
		Query().
		Where(car.ID(id)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	// Direct conversion from Ent model to domain entity
	return &entity.Car{
		ID:        carDB.ID,
		TenantID:  carDB.TenantID,
		Model:     carDB.Model,
		CreatedAt: carDB.CreatedAt,
		UpdatedAt: carDB.UpdatedAt,
	}, nil
}

// GetByIDWithTenant retrieves a car by its ID along with its tenant information
func (r *carRepository) GetByIDWithTenant(ctx context.Context, id string) (*entity.Car, error) {
	carDB, err := r.client.Car.
		Query().
		Where(car.ID(id)).
		WithTenant().
		Only(ctx)
	if err != nil {
		return nil, err
	}

	// Convert Ent model to domain entity with tenant information
	opts := repository.CarLoadOptions{WithTenant: true}
	return r.entToDomain(carDB, opts), nil
}

// Update updates an existing car
func (r *carRepository) Update(ctx context.Context, car *entity.Car) error {
	// Update the UpdatedAt field to the current time
	car.UpdatedAt = time.Now()

	_, err := r.client.Car.
		UpdateOneID(car.ID).
		SetTenantID(car.TenantID).
		SetModel(car.Model).
		SetUpdatedAt(car.UpdatedAt).
		Save(ctx)
	return err
}

// UpdateInTx updates an existing car within a transaction
func (r *carRepository) UpdateInTx(ctx context.Context, tx *entgen.Tx, car *entity.Car) error {
	// Update the UpdatedAt field to the current time
	car.UpdatedAt = time.Now()

	_, err := tx.Car.
		UpdateOneID(car.ID).
		SetTenantID(car.TenantID).
		SetModel(car.Model).
		SetUpdatedAt(car.UpdatedAt).
		Save(ctx)
	return err
}

// Delete removes a car by its ID
func (r *carRepository) Delete(ctx context.Context, id string) error {
	err := r.client.Car.
		DeleteOneID(id).
		Exec(ctx)
		// Make the delete operation idempotent by ignoring "not found" errors
		// If the record doesn't exist, DeleteOneID.Exec() will return an error
		// We want Delete to be idempotent, so we ignore "not found" errors
	if err != nil {
		// Check if it's a "not found" error by checking the error message
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "no rows in result set") {
			// Ignore "not found" errors to make the operation idempotent
			return nil
		}
		return err
	}

	return nil
}

// DeleteInTx removes a car by its ID within a transaction
func (r *carRepository) DeleteInTx(ctx context.Context, tx *entgen.Tx, id string) error {
	return tx.Car.
		DeleteOneID(id).
		Exec(ctx)
}

// ListByTenant retrieves cars by tenant ID with pagination
func (r *carRepository) ListByTenant(ctx context.Context, tenantID string, limit int, offset int) (*entity.Cars, error) {
	dbCars, err := r.client.Car.
		Query().
		Where(car.TenantID(tenantID)).
		Limit(limit).
		Offset(offset).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to query cars: %w", err)
	}

	cars := make(entity.Cars, len(dbCars))
	for i, dbCar := range dbCars {
		// Direct conversion from Ent model to domain entity
		cars[i] = &entity.Car{
			ID:        dbCar.ID,
			TenantID:  dbCar.TenantID,
			Model:     dbCar.Model,
			CreatedAt: dbCar.CreatedAt,
			UpdatedAt: dbCar.UpdatedAt,
		}
	}

	return &cars, nil
}

// ListByTenantWithOptions retrieves cars by tenant ID with pagination and load options
func (r *carRepository) ListByTenantWithOptions(ctx context.Context, tenantID string, limit int, offset int, opts ...repository.CarLoadOptions) (*entity.Cars, error) {
	query := r.client.Car.
		Query().
		Where(car.TenantID(tenantID))

	// Handle eager loading based on options
	if len(opts) > 0 {
		opt := opts[0]
		if opt.WithTenant {
			query = query.WithTenant()
		}
		if opt.WithRentals {
			query = query.WithRentals()
		}
	}

	dbCars, err := query.
		Limit(limit).
		Offset(offset).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to query cars: %w", err)
	}

	cars := make(entity.Cars, len(dbCars))
	for i, dbCar := range dbCars {
		cars[i] = r.entToDomain(dbCar, opts...)
	}

	return &cars, nil
}

// entToDomain converts an Ent car model to a domain car entity
func (r *carRepository) entToDomain(entCar *entgen.Car, opts ...repository.CarLoadOptions) *entity.Car {
	domainCar := &entity.Car{
		ID:        entCar.ID,
		TenantID:  entCar.TenantID,
		Model:     entCar.Model,
		CreatedAt: entCar.CreatedAt,
		UpdatedAt: entCar.UpdatedAt,
	}

	// Handle relationships if options are provided
	if len(opts) == 0 {
		return domainCar
	}

	opt := opts[0]
	hasTenant := opt.WithTenant && entCar.Edges.Tenant != nil
	hasRentals := opt.WithRentals && len(entCar.Edges.Rentals) > 0

	if !hasTenant && !hasRentals {
		return domainCar
	}

	domainCar.Refs = &entity.CarRefs{}

	if hasTenant {
		domainCar.Refs.Tenant = r.entTenantToDomain(entCar.Edges.Tenant)
	}

	if hasRentals {
		rentals := make(entity.Rentals, len(entCar.Edges.Rentals))
		for i, rental := range entCar.Edges.Rentals {
			rentals[i] = r.entRentalToDomain(rental)
		}
		domainCar.Refs.Rentals = rentals
	}

	return domainCar
}

// entTenantToDomain converts an Ent tenant model to a domain tenant entity
func (r *carRepository) entTenantToDomain(entTenant *entgen.Tenant) *entity.Tenant {
	return &entity.Tenant{
		ID:        entTenant.ID,
		Code:      entTenant.Code,
		CreatedAt: entTenant.CreatedAt,
		UpdatedAt: entTenant.UpdatedAt,
	}
}

// entRentalToDomain converts an Ent rental model to a domain rental entity
func (r *carRepository) entRentalToDomain(entRental *entgen.Rental) *entity.Rental {
	return &entity.Rental{
		ID:        entRental.ID,
		TenantID:  entRental.TenantID,
		CarID:     entRental.CarID,
		RenterID:  entRental.RenterID,
		StartsAt:  entRental.StartsAt,
		EndsAt:    entRental.EndsAt,
		CreatedAt: entRental.CreatedAt,
		UpdatedAt: entRental.UpdatedAt,
	}
}

// Tx returns the underlying Ent client's transaction capabilities
func (r *carRepository) Tx() *entgen.Client {
	return r.client
}

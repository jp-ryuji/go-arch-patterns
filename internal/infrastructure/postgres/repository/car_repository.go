package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	sharedrepository "github.com/jp-ryuji/go-arch-patterns/internal/domain/repository"
	"github.com/jp-ryuji/go-arch-patterns/internal/infrastructure/postgres/entgen"
	car "github.com/jp-ryuji/go-arch-patterns/internal/infrastructure/postgres/entgen/car"
	rental "github.com/jp-ryuji/go-arch-patterns/internal/rental"
)

type carRepository struct {
	client *entgen.Client
}

// NewCarRepository creates a new car repository
func NewCarRepository(client *entgen.Client) rental.CarRepository {
	return &carRepository{
		client: client,
	}
}

// Create inserts a new car into the database
func (r *carRepository) Create(ctx context.Context, car *rental.Car) error {
	_, err := r.client.Car.
		Create().
		SetID(car.ID).
		SetTenantID(car.TenantID).
		SetModel(car.Model).
		Save(ctx)
	return err
}

// CreateInTx inserts a new car into the database within a transaction
func (r *carRepository) CreateInTx(ctx context.Context, tx sharedrepository.Transaction, car *rental.Car) error {
	entTx, err := asEntTx(tx)
	if err != nil {
		return err
	}

	_, err = entTx.Car.
		Create().
		SetID(car.ID).
		SetTenantID(car.TenantID).
		SetModel(car.Model).
		Save(ctx)
	return err
}

// GetByID retrieves a car by its ID
func (r *carRepository) GetByID(ctx context.Context, id string) (*rental.Car, error) {
	carDB, err := r.client.Car.
		Query().
		Where(car.ID(id)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	// Direct conversion from Ent model to domain entity
	return &rental.Car{
		ID:        carDB.ID,
		TenantID:  carDB.TenantID,
		Model:     carDB.Model,
		CreatedAt: carDB.CreatedAt,
		UpdatedAt: carDB.UpdatedAt,
	}, nil
}

// GetByIDWithTenant retrieves a car by its ID along with its tenant information
func (r *carRepository) GetByIDWithTenant(ctx context.Context, id string) (*rental.Car, error) {
	carDB, err := r.client.Car.
		Query().
		Where(car.ID(id)).
		WithTenant().
		Only(ctx)
	if err != nil {
		return nil, err
	}

	// Convert Ent model to domain entity with tenant information
	opts := rental.CarLoadOptions{WithTenant: true}
	return r.entToDomain(carDB, opts), nil
}

// Update updates an existing car
func (r *carRepository) Update(ctx context.Context, car *rental.Car) error {
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
func (r *carRepository) UpdateInTx(ctx context.Context, tx sharedrepository.Transaction, car *rental.Car) error {
	car.UpdatedAt = time.Now()

	entTx, err := asEntTx(tx)
	if err != nil {
		return err
	}

	_, err = entTx.Car.
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
func (r *carRepository) DeleteInTx(ctx context.Context, tx sharedrepository.Transaction, id string) error {
	entTx, err := asEntTx(tx)
	if err != nil {
		return err
	}

	return entTx.Car.
		DeleteOneID(id).
		Exec(ctx)
}

// ListByTenant retrieves cars by tenant ID with pagination
func (r *carRepository) ListByTenant(ctx context.Context, tenantID string, limit int, offset int) ([]*rental.Car, string, int32, error) {
	dbCars, err := r.client.Car.
		Query().
		Where(car.TenantID(tenantID)).
		Limit(limit).
		Offset(offset).
		All(ctx)
	if err != nil {
		return nil, "", 0, fmt.Errorf("failed to query cars: %w", err)
	}

	cars := make([]*rental.Car, len(dbCars))
	for i, dbCar := range dbCars {
		// Direct conversion from Ent model to domain entity
		cars[i] = &rental.Car{
			ID:        dbCar.ID,
			TenantID:  dbCar.TenantID,
			Model:     dbCar.Model,
			CreatedAt: dbCar.CreatedAt,
			UpdatedAt: dbCar.UpdatedAt,
		}
	}

	// For now, we'll use empty nextPageToken and totalCount
	// TODO: Implement proper pagination
	// Using a simple conversion for now, in a real implementation you would get the count from the database
	count := int32(len(cars)) // #nosec G115
	return cars, "", count, nil
}

// ListByTenantWithOptions retrieves cars by tenant ID with pagination and load options
func (r *carRepository) ListByTenantWithOptions(ctx context.Context, tenantID string, limit int, offset int, opts ...rental.CarLoadOptions) ([]*rental.Car, string, int32, error) {
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
		return nil, "", 0, fmt.Errorf("failed to query cars: %w", err)
	}

	cars := make([]*rental.Car, len(dbCars))
	for i, dbCar := range dbCars {
		cars[i] = r.entToDomain(dbCar, opts...)
	}

	// For now, we'll use empty nextPageToken and totalCount
	// TODO: Implement proper pagination
	// Using a simple conversion for now, in a real implementation you would get the count from the database
	count := int32(len(cars)) // #nosec G115
	return cars, "", count, nil
}

// entToDomain converts an Ent car model to a domain car entity
func (r *carRepository) entToDomain(entCar *entgen.Car, _ ...rental.CarLoadOptions) *rental.Car {
	domainCar := &rental.Car{
		ID:        entCar.ID,
		TenantID:  entCar.TenantID,
		Model:     entCar.Model,
		CreatedAt: entCar.CreatedAt,
		UpdatedAt: entCar.UpdatedAt,
	}

	return domainCar
}

// Tx returns the underlying Ent client's transaction capabilities
func (r *carRepository) Tx() *entgen.Client {
	return r.client
}

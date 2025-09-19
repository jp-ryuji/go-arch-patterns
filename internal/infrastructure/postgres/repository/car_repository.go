package repository

import (
	"context"
	"strings"
	"time"

	"github.com/jp-ryuji/go-arch-patterns/internal/domain/entity"
	"github.com/jp-ryuji/go-arch-patterns/internal/domain/repository"
	"github.com/jp-ryuji/go-arch-patterns/internal/infrastructure/postgres/dbmodel"
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

	// Convert Ent model to dbmodel and then to domain model
	dbModel := &dbmodel.Car{
		ID:        carDB.ID,
		TenantID:  carDB.TenantID,
		Model:     carDB.Model,
		CreatedAt: carDB.CreatedAt,
		UpdatedAt: carDB.UpdatedAt,
	}
	return dbModel.ToDomain(), nil
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

	// Convert Ent model to dbmodel
	dbModel := &dbmodel.Car{
		ID:        carDB.ID,
		TenantID:  carDB.TenantID,
		Model:     carDB.Model,
		CreatedAt: carDB.CreatedAt,
		UpdatedAt: carDB.UpdatedAt,
	}

	// Load the tenant information if available
	if carDB.Edges.Tenant != nil {
		dbModel.Tenant = dbmodel.Tenant{
			ID:        carDB.Edges.Tenant.ID,
			Code:      carDB.Edges.Tenant.Code,
			CreatedAt: carDB.Edges.Tenant.CreatedAt,
			UpdatedAt: carDB.Edges.Tenant.UpdatedAt,
		}
		return dbModel.ToDomain(dbmodel.CarLoadOptions{WithTenant: true}), nil
	}

	return dbModel.ToDomain(), nil
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
	err := tx.Car.
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

// Tx returns the underlying Ent client's transaction capabilities
func (r *carRepository) Tx() *entgen.Client {
	return r.client
}

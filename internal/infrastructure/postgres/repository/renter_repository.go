package repository

import (
	"context"
	"strings"
	"time"

	"github.com/jp-ryuji/go-arch-patterns/internal/domain/entity"
	"github.com/jp-ryuji/go-arch-patterns/internal/domain/repository"
	"github.com/jp-ryuji/go-arch-patterns/internal/infrastructure/postgres/dbmodel"
	"github.com/jp-ryuji/go-arch-patterns/internal/infrastructure/postgres/entgen"
	renter "github.com/jp-ryuji/go-arch-patterns/internal/infrastructure/postgres/entgen/renter"
)

type renterRepository struct {
	client *entgen.Client
}

// NewRenterRepository creates a new renter repository
func NewRenterRepository(client *entgen.Client) repository.RenterRepository {
	return &renterRepository{
		client: client,
	}
}

// Create inserts a new renter into the database
func (r *renterRepository) Create(ctx context.Context, renter *entity.Renter) error {
	_, err := r.client.Renter.
		Create().
		SetID(renter.ID).
		SetTenantID(renter.TenantID).
		SetType(string(renter.Type)).
		SetCreatedAt(renter.CreatedAt).
		SetUpdatedAt(renter.UpdatedAt).
		Save(ctx)
	return err
}

// GetByID retrieves a renter by its ID
func (r *renterRepository) GetByID(ctx context.Context, id string) (*entity.Renter, error) {
	renterDB, err := r.client.Renter.
		Query().
		Where(renter.ID(id)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	// Convert Ent model to dbmodel and then to domain model
	dbModel := &dbmodel.Renter{
		ID:        renterDB.ID,
		TenantID:  renterDB.TenantID,
		Type:      renterDB.Type,
		CreatedAt: renterDB.CreatedAt,
		UpdatedAt: renterDB.UpdatedAt,
	}
	return dbModel.ToDomain(), nil
}

// Update updates an existing renter
func (r *renterRepository) Update(ctx context.Context, renter *entity.Renter) error {
	// Update the UpdatedAt field to the current time
	renter.UpdatedAt = time.Now()

	_, err := r.client.Renter.
		UpdateOneID(renter.ID).
		SetTenantID(renter.TenantID).
		SetType(string(renter.Type)).
		SetUpdatedAt(renter.UpdatedAt).
		Save(ctx)
	return err
}

// Delete removes a renter by its ID
func (r *renterRepository) Delete(ctx context.Context, id string) error {
	err := r.client.Renter.
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

package repository

import (
	"context"
	"strings"
	"time"

	"github.com/jp-ryuji/go-arch-patterns/internal/domain/entity"
	"github.com/jp-ryuji/go-arch-patterns/internal/domain/repository"
	"github.com/jp-ryuji/go-arch-patterns/internal/infrastructure/postgres/dbmodel"
	"github.com/jp-ryuji/go-arch-patterns/internal/infrastructure/postgres/entgen"
	company "github.com/jp-ryuji/go-arch-patterns/internal/infrastructure/postgres/entgen/company"
)

type companyRepository struct {
	client *entgen.Client
}

// NewCompanyRepository creates a new company repository
func NewCompanyRepository(client *entgen.Client) repository.CompanyRepository {
	return &companyRepository{
		client: client,
	}
}

// Create inserts a new company into the database
func (r *companyRepository) Create(ctx context.Context, company *entity.Company) error {
	_, err := r.client.Company.
		Create().
		SetID(company.ID).
		SetRenterID(company.RenterID).
		SetTenantID(company.TenantID).
		SetName(company.Name).
		SetCompanySize(company.CompanySize.String()).
		Save(ctx)
	return err
}

// GetByID retrieves a company by its ID
func (r *companyRepository) GetByID(ctx context.Context, id string) (*entity.Company, error) {
	companyDB, err := r.client.Company.
		Query().
		Where(company.RenterIDEQ(id)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	// Convert Ent model to dbmodel and then to domain model
	dbModel := &dbmodel.Company{
		ID:          companyDB.ID,
		RenterID:    companyDB.RenterID,
		TenantID:    companyDB.TenantID,
		Name:        companyDB.Name,
		CompanySize: companyDB.CompanySize,
		CreatedAt:   companyDB.CreatedAt,
		UpdatedAt:   companyDB.UpdatedAt,
	}
	return dbModel.ToDomain(), nil
}

// Update updates an existing company
func (r *companyRepository) Update(ctx context.Context, comp *entity.Company) error {
	// Update the UpdatedAt field to the current time
	comp.UpdatedAt = time.Now()

	_, err := r.client.Company.
		Update().
		Where(company.RenterIDEQ(comp.RenterID)).
		SetTenantID(comp.TenantID).
		SetName(comp.Name).
		SetCompanySize(comp.CompanySize.String()).
		SetUpdatedAt(comp.UpdatedAt).
		Save(ctx)
	return err
}

// Delete removes a company by its ID
func (r *companyRepository) Delete(ctx context.Context, id string) error {
	affected, err := r.client.Company.
		Delete().
		Where(company.RenterIDEQ(id)).
		Exec(ctx)
		// Make the delete operation idempotent by ignoring "not found" errors
		// If the record doesn't exist, Delete.Exec() will return an error
		// We want Delete to be idempotent, so we ignore "not found" errors
	if err != nil {
		// Check if it's a "not found" error by checking the error message
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "no rows in result set") {
			// Ignore "not found" errors to make the operation idempotent
			return nil
		}
		return err
	}

	// If no rows were affected, it means the record didn't exist
	// This is also an idempotent case, so we don't return an error
	_ = affected

	return nil
}

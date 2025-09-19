package dbmodel

import (
	"time"

	"github.com/aarondl/null/v9"
	"github.com/jp-ryuji/go-arch-patterns/internal/domain/entity"
	"github.com/jp-ryuji/go-arch-patterns/internal/domain/value"
)

// Individual represents the database model for Individual
type Individual struct {
	ID        string
	RenterID  string
	TenantID  string
	Email     string
	FirstName string
	LastName  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// IndividualLoadOptions specifies which associations to load
type IndividualLoadOptions struct{}

// ToDomain converts Individual to domain model with specified associations
func (i *Individual) ToDomain(opts ...IndividualLoadOptions) (*entity.Individual, error) {
	emailVO, err := value.NewEmail(i.Email)
	if err != nil {
		return nil, err
	}

	var firstName, lastName null.String
	if i.FirstName != "" {
		firstName = null.StringFrom(i.FirstName)
	}
	if i.LastName != "" {
		lastName = null.StringFrom(i.LastName)
	}

	individual := &entity.Individual{
		ID:        i.ID,
		RenterID:  i.RenterID,
		TenantID:  i.TenantID,
		Email:     *emailVO,
		FirstName: firstName,
		LastName:  lastName,
		CreatedAt: i.CreatedAt,
		UpdatedAt: i.UpdatedAt,
		Refs:      nil,
	}

	// If no options provided, return basic individual
	if len(opts) == 0 {
		return individual, nil
	}

	// Create Refs if any associations need to be loaded
	individual.Refs = &entity.IndividualRefs{}

	return individual, nil
}

// FromDomainIndividual converts domain model to Individual
func FromDomainIndividual(individual *entity.Individual) *Individual {
	firstName := ""
	if individual.FirstName.Valid {
		firstName = individual.FirstName.String
	}

	lastName := ""
	if individual.LastName.Valid {
		lastName = individual.LastName.String
	}

	return &Individual{
		ID:        individual.ID,
		RenterID:  individual.RenterID,
		TenantID:  individual.TenantID,
		Email:     individual.Email.String(),
		FirstName: firstName,
		LastName:  lastName,
		CreatedAt: individual.CreatedAt,
		UpdatedAt: individual.UpdatedAt,
	}
}

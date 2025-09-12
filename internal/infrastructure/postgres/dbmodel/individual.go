package dbmodel

import (
	"time"

	"github.com/aarondl/null/v9"
	"github.com/jp-ryuji/go-sample/internal/domain/model"
	"github.com/jp-ryuji/go-sample/internal/domain/model/value"
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

	Renters []Renter
}

// IndividualLoadOptions specifies which associations to load
type IndividualLoadOptions struct {
	WithRenters bool
}

// ToDomain converts Individual to domain model with specified associations
func (i *Individual) ToDomain(opts ...IndividualLoadOptions) (*model.Individual, error) {
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

	individual := &model.Individual{
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

	option := opts[0]

	// Only create Refs if renters need to be loaded
	if option.WithRenters && len(i.Renters) > 0 {
		renters := make(model.Renters, len(i.Renters))
		for j, renter := range i.Renters {
			renters[j] = renter.ToDomain()
		}

		individual.Refs = &model.IndividualRefs{
			Renters: renters,
		}
	}

	return individual, nil
}

// FromDomain converts domain model to Individual
func FromDomainIndividual(individual *model.Individual) *Individual {
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

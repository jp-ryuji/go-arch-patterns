package model

import (
	"time"

	"github.com/aarondl/null/v9"
	"github.com/jp-ryuji/go-sample/internal/domain/model/value"
	"github.com/jp-ryuji/go-sample/internal/pkg/id"
)

// IndividualRefs holds references to related entities for Individual
type IndividualRefs struct {
	Tenant *Tenant
}

type Individual struct {
	ID        string
	TenantID  string
	Email     *value.Email
	FirstName null.String
	LastName  null.String
	CreatedAt time.Time
	UpdatedAt time.Time

	Refs *IndividualRefs
}

// Renter interface implementation for Individual
func (i *Individual) GetID() string {
	return i.ID
}

func (i *Individual) GetEntityType() RenterEntity {
	return RenterEntityIndividual
}

// End of Renter interface implementation

type Individuals []*Individual

func NewIndividual(
	tenantID string,
	email string,
	firstName null.String,
	lastName null.String,
	t time.Time,
) (*Individual, error) {
	emailVO, err := value.NewEmail(email)
	if err != nil {
		return nil, err
	}

	return &Individual{
		ID:        id.New(),
		TenantID:  tenantID,
		Email:     emailVO,
		FirstName: firstName,
		LastName:  lastName,
		CreatedAt: t,
		UpdatedAt: t,

		Refs: nil,
	}, nil
}

package model

import (
	"time"

	"github.com/aarondl/null/v9"
	"github.com/google/uuid"

	"github.com/jp-ryuji/go-sample/internal/domain/model/value"
)

// Individuals is a slice of Individual
type Individuals []*Individual

// Individual represents an individual entity
type Individual struct {
	ID        string
	TenantID  string
	Email     value.Email
	FirstName null.String
	LastName  null.String
	CreatedAt time.Time
	UpdatedAt time.Time

	// References to related entities
	Refs *IndividualRefs
}

// IndividualRefs holds references to related entities
type IndividualRefs struct {
	Renters Renters
}

// NewIndividual creates a new Individual
func NewIndividual(tenantID string, email value.Email, firstName, lastName null.String, createdAt time.Time) *Individual {
	return &Individual{
		ID:        uuid.New().String(),
		TenantID:  tenantID,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}
}

// WithID creates an Individual with a specific ID (for testing)
func (i *Individual) WithID(id string) *Individual {
	i.ID = id
	return i
}

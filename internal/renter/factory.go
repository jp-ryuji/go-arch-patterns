package renter

import (
	"time"

	"github.com/aarondl/null/v9"

	"github.com/jp-ryuji/go-arch-patterns/internal/pkg/id"
)

func NewCompanyFixture() *Company {
	renter := NewRenter(id.New(), TypeCompany, time.Now())
	return NewCompany(
		renter.ID,
		id.New(),
		"Test Company",
		CompanySizeSmall,
		time.Now(),
	)
}

func NewIndividualFixture() (*Individual, error) {
	email, err := NewEmail("test@example.com")
	if err != nil {
		return nil, err
	}

	renter := NewRenter(id.New(), TypeIndividual, time.Now())
	return NewIndividual(
		renter.ID,
		id.New(),
		*email,
		null.StringFrom("First"),
		null.StringFrom("Last"),
		time.Now(),
	), nil
}

func NewRenterFixture() *Renter {
	return NewRenter(id.New(), TypeCompany, time.Now())
}

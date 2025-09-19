package factory

import (
	"time"

	"github.com/aarondl/null/v9"
	"github.com/jp-ryuji/go-arch-patterns/internal/domain/entity"
	"github.com/jp-ryuji/go-arch-patterns/internal/domain/value"
	"github.com/jp-ryuji/go-arch-patterns/internal/pkg/id"
)

// NewIndividual creates a new Individual for testing purposes
func NewIndividual() (*entity.Individual, error) {
	email, err := value.NewEmail("test@example.com")
	if err != nil {
		return nil, err
	}

	renter := entity.NewRenter(
		id.New(),
		entity.IndividualRenter,
		time.Now(),
	)
	return entity.NewIndividual(
		renter.ID,
		id.New(),
		*email,
		null.StringFrom("First"),
		null.StringFrom("Last"),
		time.Now(),
	), nil
}

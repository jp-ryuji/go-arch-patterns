package factory

import (
	"time"

	"github.com/aarondl/null/v9"
	"github.com/jp-ryuji/go-ddd/internal/domain/model"
	"github.com/jp-ryuji/go-ddd/internal/domain/model/value"
	"github.com/jp-ryuji/go-ddd/internal/pkg/id"
)

// NewIndividual creates a new Individual for testing purposes
func NewIndividual() (*model.Individual, error) {
	email, err := value.NewEmail("test@example.com")
	if err != nil {
		return nil, err
	}

	renter := model.NewRenter(
		id.New(),
		model.IndividualRenter,
		time.Now(),
	)
	return model.NewIndividual(
		renter.ID,
		id.New(),
		*email,
		null.StringFrom("First"),
		null.StringFrom("Last"),
		time.Now(),
	), nil
}

package factory

import (
	"time"

	"github.com/aarondl/null/v9"
	"github.com/jp-ryuji/go-sample/internal/domain/model"
	"github.com/jp-ryuji/go-sample/internal/pkg/id"
)

// NewIndividual creates a new Individual for testing purposes
func NewIndividual() (*model.Individual, error) {
	return model.NewIndividual(
		id.New(),
		"test@example.com",
		null.StringFrom("First"),
		null.StringFrom("Last"),
		time.Now(),
	)
}

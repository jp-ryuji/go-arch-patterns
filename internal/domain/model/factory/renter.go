package factory

import (
	"time"

	"github.com/jp-ryuji/go-sample/internal/domain/model"
	"github.com/jp-ryuji/go-sample/internal/pkg/id"
)

// NewRenter creates a new Renter for testing purposes
func NewRenter() *model.Renter {
	return model.NewRenter(
		id.New(),
		model.CompanyRenter, // Default to company type
		time.Now(),
	)
}

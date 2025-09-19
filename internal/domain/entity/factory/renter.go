package factory

import (
	"time"

	"github.com/jp-ryuji/go-arch-patterns/internal/domain/entity"
	"github.com/jp-ryuji/go-arch-patterns/internal/pkg/id"
)

// NewRenter creates a new Renter for testing purposes
func NewRenter() *entity.Renter {
	return entity.NewRenter(
		id.New(),
		entity.CompanyRenter, // Default to company type
		time.Now(),
	)
}

package factory

import (
	"time"

	"github.com/jp-ryuji/go-arch-patterns/internal/domain/entity"
	"github.com/jp-ryuji/go-arch-patterns/internal/pkg/id"
)

// NewCompany creates a new Company for testing purposes
func NewCompany() *entity.Company {
	renter := entity.NewRenter(
		id.New(),
		entity.CompanyRenter,
		time.Now(),
	)
	return entity.NewCompany(
		renter.ID,
		id.New(),
		"Test Company",
		entity.CompanySizeSmall,
		time.Now(),
	)
}

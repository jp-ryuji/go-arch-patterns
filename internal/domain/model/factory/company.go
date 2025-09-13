package factory

import (
	"time"

	"github.com/jp-ryuji/go-ddd/internal/domain/model"
	"github.com/jp-ryuji/go-ddd/internal/pkg/id"
)

// NewCompany creates a new Company for testing purposes
func NewCompany() *model.Company {
	renter := model.NewRenter(
		id.New(),
		model.CompanyRenter,
		time.Now(),
	)
	return model.NewCompany(
		renter.ID,
		id.New(),
		"Test Company",
		model.CompanySizeSmall,
		time.Now(),
	)
}

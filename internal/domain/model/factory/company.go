package factory

import (
	"time"

	"github.com/jp-ryuji/go-sample/internal/domain/model"
	"github.com/jp-ryuji/go-sample/internal/pkg/id"
)

// NewCompany creates a new Company for testing purposes
func NewCompany() *model.Company {
	return model.NewCompany(
		id.New(),
		"Test Company",
		model.CompanySizeSmall,
		time.Now(),
	)
}

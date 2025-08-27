package main

import (
	"gorm.io/gen"

	"github.com/jp-ryuji/go-sample/internal/infrastructure/postgres/dbmodel"
)

func main() {
	// Initialize the generator with configuration
	g := gen.NewGenerator(gen.Config{
		OutPath: "./internal/infrastructure/postgres/query",
		Mode:    gen.WithDefaultQuery,
	})

	// Apply basic CRUD interfaces to all generated models
	g.ApplyBasic(
		&dbmodel.Car{},
		&dbmodel.Company{},
		&dbmodel.Individual{},
		&dbmodel.Option{},
		&dbmodel.Rental{},
		&dbmodel.RentalOption{},
		&dbmodel.Renter{},
		&dbmodel.Tenant{},
	)

	// Execute the generator
	g.Execute()
}

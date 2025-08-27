package main

import (
	"gorm.io/gen"
)

func main() {
	// Initialize the generator with configuration
	g := gen.NewGenerator(gen.Config{
		OutPath: "./query",
		Mode:    gen.WithDefaultQuery,
	})

	// TODO: Configure database connection for GORM Gen
	// g.UseDB(getDB())

	// Apply basic CRUD interfaces to all generated models
	// g.ApplyBasic(
	// Generate queries for all models
	// g.GenerateAllTable()...,
	// )

	// Execute the generator
	g.Execute()
}

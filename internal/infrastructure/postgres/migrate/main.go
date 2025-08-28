package main

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/jp-ryuji/go-sample/internal/infrastructure/postgres/dbmodel"
)

func main() {
	// Get database connection parameters from environment variables
	host := getEnv("DB_HOST", "localhost")
	user := getEnv("DB_USER", "user")
	password := getEnv("DB_PASSWORD", "password")
	dbname := getEnv("DB_NAME", "mydb")
	port := getEnv("DB_PORT_EXTERNAL", "5433") // Use external port
	sslmode := getEnv("DB_SSLMODE", "disable")

	// Create connection string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslmode)

	// Connect to database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	// Migrate all schemas
	if err := db.AutoMigrate(
		&dbmodel.Car{},
		&dbmodel.Company{},
		&dbmodel.Individual{},
		&dbmodel.Option{},
		&dbmodel.Rental{},
		&dbmodel.RentalOption{},
		&dbmodel.Renter{},
		&dbmodel.Tenant{},
	); err != nil {
		log.Fatal("failed to migrate database:", err)
	}

	// Drop foreign key constraints on renters table that were automatically created by GORM
	// These constraints are problematic for polymorphic relationships
	if err := db.Exec("ALTER TABLE renters DROP CONSTRAINT IF EXISTS fk_companies_renters").Error; err != nil {
		log.Printf("Warning: failed to drop fk_companies_renters: %v", err)
	}
	if err := db.Exec("ALTER TABLE renters DROP CONSTRAINT IF EXISTS fk_individuals_renters").Error; err != nil {
		log.Printf("Warning: failed to drop fk_individuals_renters: %v", err)
	}

	// Drop duplicate foreign key constraints
	if err := db.Exec("ALTER TABLE rentals DROP CONSTRAINT IF EXISTS fk_renters_rentals").Error; err != nil {
		log.Printf("Warning: failed to drop fk_renters_rentals: %v", err)
	}
	if err := db.Exec("ALTER TABLE cars DROP CONSTRAINT IF EXISTS fk_tenants_cars").Error; err != nil {
		log.Printf("Warning: failed to drop fk_tenants_cars: %v", err)
	}
	if err := db.Exec("ALTER TABLE rentals DROP CONSTRAINT IF EXISTS fk_rentals_car").Error; err != nil {
		log.Printf("Warning: failed to drop fk_rentals_car: %v", err)
	}

	fmt.Println("Database migration completed")
}

// getEnv returns the value of the environment variable or a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

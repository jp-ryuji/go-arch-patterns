package main

import (
	"context"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jp-ryuji/go-ddd/internal/infrastructure/postgres"
)

func main() {
	// Get database connection details from environment variables
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "user")
	password := getEnv("DB_PASSWORD", "password")
	dbname := getEnv("DB_NAME", "postgres")
	sslmode := getEnv("DB_SSLMODE", "disable")

	// Create connection string
	databaseUrl := "postgres://" + user + ":" + password + "@" + host + ":" + port + "/" + dbname + "?sslmode=" + sslmode

	// Create Ent client using our NewClient function
	client := postgres.NewClient(databaseUrl)
	defer client.Close()

	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	log.Println("Migration completed successfully")
}

// getEnv returns the value of the environment variable named by the key,
// or the fallback value if the environment variable is not set.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

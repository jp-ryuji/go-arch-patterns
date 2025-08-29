package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jp-ryuji/go-sample/internal/infrastructure/postgres/entgen"
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
	connStr := "postgres://" + user + ":" + password + "@" + host + ":" + port + "/" + dbname + "?sslmode=" + sslmode

	// Open database connection with pgx driver
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer db.Close()

	// Create Ent driver with the database connection
	drv := entsql.OpenDB(dialect.Postgres, db)

	// Create Ent client with the driver
	client := entgen.NewClient(entgen.Driver(drv))

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

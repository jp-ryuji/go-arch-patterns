package postgres

import (
	"context"
	"database/sql"
	"log"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jp-ryuji/go-sample/internal/infrastructure/postgres/entgen"
)

// NewClient creates a new Ent client with pgx driver
func NewClient(databaseUrl string) *entgen.Client {
	log.Printf("Connecting to database with connection string: %s", databaseUrl)

	// Create database connection with pgx driver
	db, err := sql.Open("pgx", databaseUrl)
	if err != nil {
		log.Printf("Failed to open SQL DB: %v", err)
		panic(err)
	}

	// Ping to verify connection
	log.Printf("Pinging database...")
	if err := db.PingContext(context.Background()); err != nil {
		log.Printf("Failed to ping database: %v", err)
		panic(err)
	}

	// Create Ent driver with the database connection
	drv := entsql.OpenDB(dialect.Postgres, db)

	// Create Ent client with the driver
	entClient := entgen.NewClient(entgen.Driver(drv))

	log.Printf("Successfully connected to database")

	return entClient
}

package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jp-ryuji/go-sample/internal/infrastructure/postgres/entgen"
)

// Client wraps the Ent client
type Client struct {
	EntClient *entgen.Client
}

func NewClient(
	host,
	user,
	password,
	database string,
	logEnable bool,
) *Client {
	sslmode := "disable"
	// FIXME: Enable sslmode in production
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		user,
		password,
		host,
		database,
		sslmode)

	log.Printf("Connecting to database with connection string: %s", connStr)

	// Create database connection with pgx driver
	db, err := sql.Open("pgx", connStr)
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

	return &Client{
		EntClient: entClient,
	}
}

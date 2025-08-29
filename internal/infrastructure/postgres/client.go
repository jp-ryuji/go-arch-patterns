package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"entgo.io/ent/dialect"
	"github.com/jp-ryuji/go-sample/internal/infrastructure/postgres/entgen"
	_ "github.com/lib/pq"
)

type Client struct {
	EntClient *entgen.Client
	DB        *sql.DB
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

	// Create Ent client with "postgres" dialect
	entClient, err := entgen.Open(dialect.Postgres, connStr)
	if err != nil {
		log.Printf("Failed to create Ent client: %v", err)
		panic(err)
	}

	// Get the underlying sql.DB from ent for compatibility
	sqlDB, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Failed to open SQL DB: %v", err)
		panic(err)
	}

	// Ping to verify connection
	log.Printf("Pinging database...")
	if err := sqlDB.PingContext(context.Background()); err != nil {
		log.Printf("Failed to ping database: %v", err)
		panic(err)
	}

	log.Printf("Successfully connected to database")

	return &Client{
		EntClient: entClient,
		DB:        sqlDB,
	}
}

package postgres

import (
	"context"
	"database/sql"
	"log"
	"os"
	"strconv"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jp-ryuji/go-ddd/internal/infrastructure/postgres/entgen"
)

// Connection pool settings with defaults
// Reference: https://entgo.io/docs/sql-integration/
const (
	defaultMaxOpenConns    = 25
	defaultMaxIdleConns    = 25
	defaultConnMaxLifetime = 300 * time.Second // 5 minutes (max connection * seconds)
)

// getConnectionPoolSettings reads connection pool settings from environment variables
// with fallback to default values
func getConnectionPoolSettings() (maxOpenConns, maxIdleConns int, maxLifetime time.Duration) {
	// Max open connections
	maxOpenConns = defaultMaxOpenConns
	if envMaxOpenConns := os.Getenv("DB_MAX_OPEN_CONNS"); envMaxOpenConns != "" {
		if n, err := strconv.Atoi(envMaxOpenConns); err == nil && n > 0 {
			maxOpenConns = n
		}
	}

	// Max idle connections
	maxIdleConns = defaultMaxIdleConns
	if envMaxIdleConns := os.Getenv("DB_MAX_IDLE_CONNS"); envMaxIdleConns != "" {
		if n, err := strconv.Atoi(envMaxIdleConns); err == nil && n > 0 {
			maxIdleConns = n
		}
	}

	// Connection max lifetime
	maxLifetime = defaultConnMaxLifetime
	if envMaxLifetime := os.Getenv("DB_CONN_MAX_LIFETIME"); envMaxLifetime != "" {
		// Try to parse as duration (e.g., "300s", "5m")
		if d, err := time.ParseDuration(envMaxLifetime); err == nil && d > 0 {
			maxLifetime = d
		}
	}

	return maxOpenConns, maxIdleConns, maxLifetime
}

// NewClient creates a new Ent client with pgx driver
func NewClient(databaseUrl string) *entgen.Client {
	log.Printf("Connecting to database with connection string: %s", databaseUrl)

	// Create database connection with pgx driver
	db, err := sql.Open("pgx", databaseUrl)
	if err != nil {
		log.Printf("Failed to open SQL DB: %v", err)
		panic(err)
	}

	// Configure connection pool settings
	maxOpenConns, maxIdleConns, maxLifetime := getConnectionPoolSettings()
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(maxLifetime)

	log.Printf("Database connection pool settings: MaxOpenConns=%d, MaxIdleConns=%d, ConnMaxLifetime=%v",
		maxOpenConns, maxIdleConns, maxLifetime)

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

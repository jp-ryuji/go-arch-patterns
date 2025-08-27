package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // pgx driver for database/sql
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/jp-ryuji/go-sample/internal/infrastructure/postgres/query"
)

const (
	maxOpenConns = 25
	maxIdleConns = 25
	maxLifeTime  = 300 * time.Second // max connection * seconds
)

type Client struct {
	DB *gorm.DB
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

	sqlDB, err := sql.Open("pgx", connStr)
	if err != nil {
		panic(err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetConnMaxLifetime(maxLifeTime)

	// Ping to verify connection
	if err := sqlDB.PingContext(context.Background()); err != nil {
		panic(err)
	}

	// Configure GORM logger
	gormLogger := logger.Default.LogMode(logger.Silent)
	if logEnable {
		gormLogger = logger.Default.LogMode(logger.Info)
	}

	// Create GORM DB instance
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		panic(err)
	}

	// Initialize the query package with the database connection
	// This makes the type-safe generated code globally accessible through
	// query.Car, query.Company, etc. This is a convenience pattern that
	// allows direct access to type-safe queries from anywhere in the application.
	query.SetDefault(gormDB)

	return &Client{
		DB: gormDB,
	}
}

package database

import (
	"context"
	"log"

	"github.com/jp-ryuji/go-sample/internal/ent"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// NewClient creates a new ent client
func NewClient() *ent.Client {
	// For now, we'll use an in-memory SQLite database for simplicity
	// In a production environment, you would use a real database connection
	client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}

	// Run the auto migration tool
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	return client
}

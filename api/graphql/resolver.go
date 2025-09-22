package graphql

import (
	"github.com/jp-ryuji/go-arch-patterns/internal/infrastructure/postgres/entgen"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	client *entgen.Client
}

// NewResolver creates a new GraphQL resolver
func NewResolver(client *entgen.Client) *Resolver {
	return &Resolver{
		client: client,
	}
}

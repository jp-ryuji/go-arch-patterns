package graphql

import (
	"github.com/jp-ryuji/go-sample/internal/ent"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	client *ent.Client
}

// NewResolver creates a new Resolver with the given ent client
func NewResolver(client *ent.Client) *Resolver {
	return &Resolver{client: client}
}

package graphql

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/jp-ryuji/go-arch-patterns/api/graphql"
	"github.com/jp-ryuji/go-arch-patterns/api/graphql/generated"
	"github.com/jp-ryuji/go-arch-patterns/internal/infrastructure/postgres/entgen"
)

// NewHandler creates a new GraphQL HTTP handler
func NewHandler(client *entgen.Client) http.Handler {
	resolver := graphql.NewResolver(client)
	config := generated.Config{Resolvers: resolver}

	return handler.NewDefaultServer(generated.NewExecutableSchema(config))
}

// NewPlaygroundHandler creates a GraphQL playground handler for development
func NewPlaygroundHandler() http.Handler {
	return playground.Handler("GraphQL playground", "/graphql")
}

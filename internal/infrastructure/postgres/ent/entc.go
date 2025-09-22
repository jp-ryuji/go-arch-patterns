//go:build ignore

package main

import (
	"log"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

func main() {
	ex, err := entgql.NewExtension(
		// Tell Ent to generate a GraphQL schema for
		// the Ent schema in a file named ent.graphql.
		entgql.WithSchemaGenerator(),
		entgql.WithSchemaPath("../entgen/ent.graphql"),
		entgql.WithConfigPath("../../../../api/graphql/gqlgen.yml"),
	)
	if err != nil {
		log.Fatalf("creating entgql extension: %v", err)
	}
	opts := []entc.Option{
		entc.Extensions(ex),
		entc.TemplateDir("./custom.tmpl"),
	}
	if err := entc.Generate("./schema", &gen.Config{
		Target:  "../entgen",
		Package: "github.com/jp-ryuji/go-arch-patterns/internal/infrastructure/postgres/entgen",
	}, opts...); err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}
}

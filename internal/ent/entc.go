//go:build ignore

package main

import (
	"log"
	"os"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

func main() {
	// Create the GraphQL directory if it doesn't exist
	err := os.MkdirAll("../../graphql", 0o755)
	if err != nil {
		log.Fatalf("creating graphql directory: %v", err)
	}

	ex, err := entgql.NewExtension(
		entgql.WithSchemaGenerator(),
		entgql.WithSchemaPath("../../schema/graphql/ent.graphql"),
	)
	if err != nil {
		log.Fatalf("creating entgql extension: %v", err)
	}
	opts := []entc.Option{
		entc.Extensions(ex),
	}
	err = entc.Generate("../schema", &gen.Config{
		Target: "./",
		Features: []gen.Feature{
			gen.FeatureUpsert,
			gen.FeatureExecQuery,
		},
	}, opts...)
	if err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}
}

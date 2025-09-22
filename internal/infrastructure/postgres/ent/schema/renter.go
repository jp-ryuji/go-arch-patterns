package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Renter holds the schema definition for the Renter entity.
type Renter struct {
	ent.Schema
}

// Fields of the Renter.
func (Renter) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			MaxLen(36).
			NotEmpty().
			Annotations(entgql.OrderField("ID")),
		field.String("tenant_id").
			MaxLen(36).
			NotEmpty().
			Annotations(entgql.OrderField("TENANT_ID")),
		field.String("type").
			MaxLen(20).
			NotEmpty().
			Annotations(entgql.OrderField("TYPE")),
		field.Time("created_at").
			Optional().
			Annotations(entgql.OrderField("CREATED_AT")),
		field.Time("updated_at").
			Optional().
			Annotations(entgql.OrderField("UPDATED_AT")),
		field.Time("deleted_at").
			Optional().
			Nillable().
			Annotations(entgql.Skip()), // Skip deleted_at in GraphQL
	}
}

// Edges of the Renter.
func (Renter) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tenant", Tenant.Type).
			Ref("renters").
			Field("tenant_id").
			Required().
			Unique(),
		edge.To("rentals", Rental.Type),
		edge.To("company", Company.Type).
			Unique(),
		edge.To("individual", Individual.Type).
			Unique(),
	}
}

// Indexes of the Renter.
func (Renter) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("deleted_at"),
		index.Fields("tenant_id"),
	}
}

// Annotations of the Renter.
func (Renter) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.RelayConnection(),
		entgql.QueryField(),
		entgql.Mutations(entgql.MutationCreate(), entgql.MutationUpdate()),
		// Add soft delete support
		entgql.MultiOrder(),
	}
}

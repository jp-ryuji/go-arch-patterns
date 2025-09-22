package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Tenant holds the schema definition for the Tenant entity.
type Tenant struct {
	ent.Schema
}

// Fields of the Tenant.
func (Tenant) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			MaxLen(36).
			NotEmpty().
			Annotations(entgql.OrderField("ID")),
		field.String("code").
			MaxLen(50).
			NotEmpty().
			Annotations(entgql.OrderField("CODE")),
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

// Edges of the Tenant.
func (Tenant) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("cars", Car.Type),
		edge.To("companies", Company.Type),
		edge.To("individuals", Individual.Type),
		edge.To("options", CarOption.Type),
		edge.To("rental_options", RentalOption.Type),
		edge.To("rentals", Rental.Type),
		edge.To("renters", Renter.Type),
	}
}

// Indexes of the Tenant.
func (Tenant) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("code").
			Unique(),
		index.Fields("deleted_at"),
	}
}

// Annotations of the Tenant.
func (Tenant) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.RelayConnection(),
		entgql.QueryField(),
		entgql.Mutations(entgql.MutationCreate(), entgql.MutationUpdate()),
		// Add soft delete support
		entgql.MultiOrder(),
	}
}

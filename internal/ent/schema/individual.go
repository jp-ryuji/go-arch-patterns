package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Individual holds the schema definition for the Individual entity.
type Individual struct {
	ent.Schema
}

// Fields of the Individual.
func (Individual) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique().
			NotEmpty().
			Annotations(
				entgql.OrderField("ID"),
			),
		field.String("email").
			Unique().
			Annotations(
				entgql.OrderField("EMAIL"),
			),
		field.String("first_name").
			Optional().
			Annotations(
				entgql.OrderField("FIRST_NAME"),
			),
		field.String("last_name").
			Optional().
			Annotations(
				entgql.OrderField("LAST_NAME"),
			),
		field.String("tenant_id").
			Annotations(
				entgql.OrderField("TENANT_ID"),
			),
	}
}

// Edges of the Individual.
func (Individual) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("tenant", Tenant.Type).
			Field("tenant_id").
			Unique().
			Required(),
		edge.From("renters", Renter.Type).
			Ref("individual"),
	}
}

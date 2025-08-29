package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Renter holds the schema definition for the Renter entity.
type Renter struct {
	ent.Schema
}

// Fields of the Renter.
func (Renter) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique().
			NotEmpty().
			Annotations(
				entgql.OrderField("ID"),
			),
		field.String("renter_entity_id").
			Optional().
			Annotations(
				entgql.OrderField("RENTER_ENTITY_ID"),
			),
		field.String("renter_entity_type").
			Optional().
			Annotations(
				entgql.OrderField("RENTER_ENTITY_TYPE"),
			),
		field.String("tenant_id").
			Annotations(
				entgql.OrderField("TENANT_ID"),
			),
	}
}

// Edges of the Renter.
func (Renter) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("tenant", Tenant.Type).
			Field("tenant_id").
			Unique().
			Required(),
		edge.From("rentals", Rental.Type).
			Ref("renter"),
		edge.To("company", Company.Type).
			Field("renter_entity_id").
			Unique(),
		edge.To("individual", Individual.Type).
			Field("renter_entity_id").
			Unique(),
	}
}

package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// RentalOption holds the schema definition for the RentalOption entity.
type RentalOption struct {
	ent.Schema
}

// Fields of the RentalOption.
func (RentalOption) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique().
			NotEmpty().
			Annotations(
				entgql.OrderField("ID"),
			),
		field.Int("count").
			Annotations(
				entgql.OrderField("COUNT"),
			),
		field.String("rental_id").
			Annotations(
				entgql.OrderField("RENTAL_ID"),
			),
		field.String("option_id").
			Annotations(
				entgql.OrderField("OPTION_ID"),
			),
		field.String("tenant_id").
			Annotations(
				entgql.OrderField("TENANT_ID"),
			),
	}
}

// Edges of the RentalOption.
func (RentalOption) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("tenant", Tenant.Type).
			Field("tenant_id").
			Unique().
			Required(),
		edge.To("rental", Rental.Type).
			Field("rental_id").
			Unique().
			Required(),
		edge.To("option", CarOption.Type).
			Field("option_id").
			Unique().
			Required(),
	}
}

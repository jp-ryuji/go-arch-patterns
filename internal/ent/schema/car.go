package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Car holds the schema definition for the Car entity.
type Car struct {
	ent.Schema
}

// Fields of the Car.
func (Car) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique().
			NotEmpty().
			Annotations(
				entgql.OrderField("ID"),
			),
		field.String("model").
			NotEmpty().
			Annotations(
				entgql.OrderField("MODEL"),
			),
		field.String("tenant_id").
			Annotations(
				entgql.OrderField("TENANT_ID"),
			),
	}
}

// Edges of the Car.
func (Car) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("tenant", Tenant.Type).
			Field("tenant_id").
			Unique().
			Required(),
		edge.From("rentals", Rental.Type).
			Ref("car"),
	}
}

package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Rental holds the schema definition for the Rental entity.
type Rental struct {
	ent.Schema
}

// Fields of the Rental.
func (Rental) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique().
			NotEmpty().
			Annotations(
				entgql.OrderField("ID"),
			),
		field.Time("starts_at").
			Annotations(
				entgql.OrderField("STARTS_AT"),
			),
		field.Time("ends_at").
			Annotations(
				entgql.OrderField("ENDS_AT"),
			),
		field.String("car_id").
			Annotations(
				entgql.OrderField("CAR_ID"),
			),
		field.String("renter_id").
			Annotations(
				entgql.OrderField("RENTER_ID"),
			),
		field.String("tenant_id").
			Annotations(
				entgql.OrderField("TENANT_ID"),
			),
	}
}

// Edges of the Rental.
func (Rental) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("tenant", Tenant.Type).
			Field("tenant_id").
			Unique().
			Required(),
		edge.To("car", Car.Type).
			Field("car_id").
			Unique().
			Required(),
		edge.To("renter", Renter.Type).
			Field("renter_id").
			Unique().
			Required(),
		edge.From("rental_options", RentalOption.Type).
			Ref("rental"),
	}
}

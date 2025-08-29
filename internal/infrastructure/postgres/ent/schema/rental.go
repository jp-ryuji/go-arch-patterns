package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Rental holds the schema definition for the Rental entity.
type Rental struct {
	ent.Schema
}

// Fields of the Rental.
func (Rental) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			MaxLen(36).
			NotEmpty(),
		field.String("tenant_id").
			MaxLen(36).
			NotEmpty(),
		field.String("car_id").
			MaxLen(36).
			NotEmpty(),
		field.String("renter_id").
			MaxLen(36).
			NotEmpty(),
		field.Time("starts_at").
			Optional(),
		field.Time("ends_at").
			Optional(),
		field.Time("created_at").
			Optional(),
		field.Time("updated_at").
			Optional(),
		field.Time("deleted_at").
			Optional().
			Nillable(),
	}
}

// Edges of the Rental.
func (Rental) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tenant", Tenant.Type).
			Ref("rentals").
			Field("tenant_id").
			Required().
			Unique(),
		edge.From("car", Car.Type).
			Ref("rentals").
			Field("car_id").
			Required().
			Unique(),
		edge.From("renter", Renter.Type).
			Ref("rentals").
			Field("renter_id").
			Required().
			Unique(),
		edge.To("rental_options", RentalOption.Type),
	}
}

// Indexes of the Rental.
func (Rental) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("car_id"),
		index.Fields("deleted_at"),
		index.Fields("renter_id"),
		index.Fields("tenant_id"),
	}
}

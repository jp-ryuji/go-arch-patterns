package schema

import (
	"entgo.io/ent"
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
			NotEmpty(),
		field.String("code").
			MaxLen(50).
			NotEmpty(),
		field.Time("created_at").
			Optional(),
		field.Time("updated_at").
			Optional(),
		field.Time("deleted_at").
			Optional().
			Nillable(),
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

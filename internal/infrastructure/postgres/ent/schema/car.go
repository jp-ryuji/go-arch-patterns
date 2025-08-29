package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Car holds the schema definition for the Car entity.
type Car struct {
	ent.Schema
}

// Fields of the Car.
func (Car) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			MaxLen(36).
			NotEmpty(),
		field.String("tenant_id").
			MaxLen(36).
			NotEmpty(),
		field.String("model").
			MaxLen(255).
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

// Edges of the Car.
func (Car) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tenant", Tenant.Type).
			Ref("cars").
			Field("tenant_id").
			Required().
			Unique(),
		edge.To("rentals", Rental.Type),
	}
}

// Indexes of the Car.
func (Car) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("tenant_id", "model").
			Unique(),
		index.Fields("deleted_at"),
		index.Fields("tenant_id"),
	}
}

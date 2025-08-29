package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// RentalOption holds the schema definition for the RentalOption entity.
type RentalOption struct {
	ent.Schema
}

// Fields of the RentalOption.
func (RentalOption) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			MaxLen(36).
			NotEmpty(),
		field.String("tenant_id").
			MaxLen(36).
			NotEmpty(),
		field.String("rental_id").
			MaxLen(36).
			NotEmpty(),
		field.String("option_id").
			MaxLen(36).
			NotEmpty(),
		field.Int("count"),
		field.Time("created_at").
			Optional(),
		field.Time("updated_at").
			Optional(),
		field.Time("deleted_at").
			Optional().
			Nillable(),
	}
}

// Edges of the RentalOption.
func (RentalOption) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tenant", Tenant.Type).
			Ref("rental_options").
			Field("tenant_id").
			Required().
			Unique(),
		edge.From("rental", Rental.Type).
			Ref("rental_options").
			Field("rental_id").
			Required().
			Unique(),
		edge.From("option", CarOption.Type).
			Ref("rental_options").
			Field("option_id").
			Required().
			Unique(),
	}
}

// Indexes of the RentalOption.
func (RentalOption) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("rental_id", "option_id").
			Unique(),
		index.Fields("deleted_at"),
		index.Fields("tenant_id"),
	}
}

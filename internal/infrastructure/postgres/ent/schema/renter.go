package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Renter holds the schema definition for the Renter entity.
type Renter struct {
	ent.Schema
}

// Fields of the Renter.
func (Renter) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			MaxLen(36).
			NotEmpty(),
		field.String("tenant_id").
			MaxLen(36).
			NotEmpty(),
		field.String("type").
			MaxLen(20).
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

// Edges of the Renter.
func (Renter) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tenant", Tenant.Type).
			Ref("renters").
			Field("tenant_id").
			Required().
			Unique(),
		edge.To("rentals", Rental.Type),
		edge.To("company", Company.Type).
			Unique(),
		edge.To("individual", Individual.Type).
			Unique(),
	}
}

// Indexes of the Renter.
func (Renter) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("deleted_at"),
		index.Fields("tenant_id"),
	}
}

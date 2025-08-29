package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// CarOption holds the schema definition for the CarOption entity.
type CarOption struct {
	ent.Schema
}

// Fields of the CarOption.
func (CarOption) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			MaxLen(36).
			NotEmpty(),
		field.String("tenant_id").
			MaxLen(36).
			NotEmpty(),
		field.String("name").
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

// Edges of the CarOption.
func (CarOption) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tenant", Tenant.Type).
			Ref("options").
			Field("tenant_id").
			Required().
			Unique(),
		edge.To("rental_options", RentalOption.Type),
	}
}

// Indexes of the CarOption.
func (CarOption) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("deleted_at"),
		index.Fields("tenant_id"),
	}
}

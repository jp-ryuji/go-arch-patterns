package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Company holds the schema definition for the Company entity.
type Company struct {
	ent.Schema
}

// Fields of the Company.
func (Company) Fields() []ent.Field {
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
		field.String("company_size").
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

// Edges of the Company.
func (Company) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tenant", Tenant.Type).
			Ref("companies").
			Field("tenant_id").
			Required().
			Unique(),
		edge.To("renters", Renter.Type),
	}
}

// Indexes of the Company.
func (Company) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("deleted_at"),
		index.Fields("tenant_id"),
	}
}

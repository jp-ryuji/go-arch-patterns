package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Individual holds the schema definition for the Individual entity.
type Individual struct {
	ent.Schema
}

// Fields of the Individual.
func (Individual) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			MaxLen(36).
			NotEmpty(),
		field.String("renter_id").
			MaxLen(36).
			NotEmpty().
			Unique().
			StructTag(`json:"renter_id"`).
			StorageKey("renter_id"),
		field.String("tenant_id").
			MaxLen(36).
			NotEmpty(),
		field.String("email").
			MaxLen(255).
			NotEmpty(),
		field.String("first_name").
			MaxLen(100).
			Optional(),
		field.String("last_name").
			MaxLen(100).
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

// Edges of the Individual.
func (Individual) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tenant", Tenant.Type).
			Ref("individuals").
			Field("tenant_id").
			Required().
			Unique(),
		edge.To("renter", Renter.Type).
			Field("renter_id").
			Required().
			Unique(),
	}
}

// Indexes of the Individual.
func (Individual) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("renter_id").
			Unique(),
		index.Fields("tenant_id", "email").
			Unique(),
		index.Fields("deleted_at"),
	}
}

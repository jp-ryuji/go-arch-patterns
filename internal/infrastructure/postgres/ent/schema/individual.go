package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
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
			NotEmpty().
			Annotations(entgql.OrderField("ID")),
		field.String("renter_id").
			MaxLen(36).
			NotEmpty().
			Unique().
			StructTag(`json:"renter_id"`).
			StorageKey("renter_id").
			Annotations(entgql.OrderField("RENTER_ID")),
		field.String("tenant_id").
			MaxLen(36).
			NotEmpty().
			Annotations(entgql.OrderField("TENANT_ID")),
		field.String("email").
			MaxLen(255).
			NotEmpty().
			Annotations(entgql.OrderField("EMAIL")),
		field.String("first_name").
			MaxLen(100).
			Optional().
			Annotations(entgql.OrderField("FIRST_NAME")),
		field.String("last_name").
			MaxLen(100).
			Optional().
			Annotations(entgql.OrderField("LAST_NAME")),
		field.Time("created_at").
			Optional().
			Annotations(entgql.OrderField("CREATED_AT")),
		field.Time("updated_at").
			Optional().
			Annotations(entgql.OrderField("UPDATED_AT")),
		field.Time("deleted_at").
			Optional().
			Nillable().
			Annotations(entgql.Skip()), // Skip deleted_at in GraphQL
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

// Annotations of the Individual.
func (Individual) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.RelayConnection(),
		entgql.QueryField(),
		entgql.Mutations(entgql.MutationCreate(), entgql.MutationUpdate()),
		// Add soft delete support
		entgql.MultiOrder(),
	}
}

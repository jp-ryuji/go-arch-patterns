package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
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
			NotEmpty().
			Annotations(entgql.OrderField("ID")),
		field.String("tenant_id").
			MaxLen(36).
			NotEmpty().
			Annotations(entgql.OrderField("TENANT_ID")),
		field.String("name").
			MaxLen(255).
			NotEmpty().
			Annotations(entgql.OrderField("NAME")),
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

// Annotations of the CarOption.
func (CarOption) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.RelayConnection(),
		entgql.QueryField(),
		entgql.Mutations(entgql.MutationCreate(), entgql.MutationUpdate()),
		// Add soft delete support
		entgql.MultiOrder(),
	}
}

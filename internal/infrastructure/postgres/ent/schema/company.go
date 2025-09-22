package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
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
		field.String("name").
			MaxLen(255).
			NotEmpty().
			Annotations(entgql.OrderField("NAME")),
		field.String("company_size").
			MaxLen(50).
			NotEmpty().
			Annotations(entgql.OrderField("COMPANY_SIZE")),
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

// Edges of the Company.
func (Company) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tenant", Tenant.Type).
			Ref("companies").
			Field("tenant_id").
			Required().
			Unique(),
		edge.To("renter", Renter.Type).
			Field("renter_id").
			Required().
			Unique(),
	}
}

// Indexes of the Company.
func (Company) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("renter_id").
			Unique(),
		index.Fields("deleted_at"),
		index.Fields("tenant_id"),
	}
}

// Annotations of the Company.
func (Company) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.RelayConnection(),
		entgql.QueryField(),
		entgql.Mutations(entgql.MutationCreate(), entgql.MutationUpdate()),
		// Add soft delete support
		entgql.MultiOrder(),
	}
}

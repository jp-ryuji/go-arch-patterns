package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Outbox struct {
	ent.Schema
}

// Fields of the Outbox.
func (Outbox) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			MaxLen(36).
			NotEmpty().
			Annotations(entgql.OrderField("ID")),
		field.String("aggregate_type").
			MaxLen(255).
			NotEmpty().
			Annotations(entgql.OrderField("AGGREGATE_TYPE")),
		field.String("aggregate_id").
			MaxLen(36).
			NotEmpty().
			Annotations(entgql.OrderField("AGGREGATE_ID")),
		field.String("event_type").
			MaxLen(255).
			NotEmpty().
			Annotations(entgql.OrderField("EVENT_TYPE")),
		field.JSON("payload", map[string]interface{}{}).
			Optional(),
		field.Time("created_at").
			Optional().
			Annotations(entgql.OrderField("CREATED_AT")),
		field.Time("processed_at").
			Optional().
			Nillable().
			Annotations(entgql.OrderField("PROCESSED_AT")),
		field.String("status").
			MaxLen(50).
			Default("pending").
			Annotations(entgql.OrderField("STATUS")),
		field.String("error_message").
			MaxLen(1000).
			Optional().
			Nillable().
			Annotations(entgql.OrderField("ERROR_MESSAGE")),
		field.Int64("version").
			Default(1).
			Annotations(entgql.OrderField("VERSION")),
		field.Time("locked_at").
			Optional().
			Nillable().
			Annotations(entgql.OrderField("LOCKED_AT")),
		field.String("locked_by").
			Optional().
			Nillable().
			Annotations(entgql.OrderField("LOCKED_BY")),
	}
}

// Indexes of the Outbox.
func (Outbox) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("status"),
		index.Fields("created_at"),
		index.Fields("processed_at"),
		index.Fields("version"),
		index.Fields("locked_at", "locked_by"),
	}
}

// Annotations of the Outbox.
func (Outbox) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.RelayConnection(),
		entgql.QueryField(),
		entgql.Mutations(entgql.MutationCreate(), entgql.MutationUpdate()),
		// Add soft delete support
		entgql.MultiOrder(),
	}
}

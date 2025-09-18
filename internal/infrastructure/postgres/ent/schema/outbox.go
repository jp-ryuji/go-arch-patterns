package schema

import (
	"entgo.io/ent"
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
			NotEmpty(),
		field.String("aggregate_type").
			MaxLen(255).
			NotEmpty(),
		field.String("aggregate_id").
			MaxLen(36).
			NotEmpty(),
		field.String("event_type").
			MaxLen(255).
			NotEmpty(),
		field.JSON("payload", map[string]interface{}{}).
			Optional(),
		field.Time("created_at").
			Optional(),
		field.Time("processed_at").
			Optional().
			Nillable(),
		field.String("status").
			MaxLen(50).
			Default("pending"),
		field.String("error_message").
			MaxLen(1000).
			Optional().
			Nillable(),
		field.Int64("version").
			Default(1),
		field.Time("locked_at").
			Optional().
			Nillable(),
		field.String("locked_by").
			Optional().
			Nillable(),
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

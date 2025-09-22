package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Rental holds the schema definition for the Rental entity.
type Rental struct {
	ent.Schema
}

// Fields of the Rental.
func (Rental) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			MaxLen(36).
			NotEmpty().
			Annotations(entgql.OrderField("ID")),
		field.String("tenant_id").
			MaxLen(36).
			NotEmpty().
			Annotations(entgql.OrderField("TENANT_ID")),
		field.String("car_id").
			MaxLen(36).
			NotEmpty().
			Annotations(entgql.OrderField("CAR_ID")),
		field.String("renter_id").
			MaxLen(36).
			NotEmpty().
			Annotations(entgql.OrderField("RENTER_ID")),
		field.Time("starts_at").
			Optional().
			Annotations(entgql.OrderField("STARTS_AT")),
		field.Time("ends_at").
			Optional().
			Annotations(entgql.OrderField("ENDS_AT")),
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

// Edges of the Rental.
func (Rental) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tenant", Tenant.Type).
			Ref("rentals").
			Field("tenant_id").
			Required().
			Unique(),
		edge.From("car", Car.Type).
			Ref("rentals").
			Field("car_id").
			Required().
			Unique(),
		edge.From("renter", Renter.Type).
			Ref("rentals").
			Field("renter_id").
			Required().
			Unique(),
		edge.To("rental_options", RentalOption.Type),
	}
}

// Indexes of the Rental.
func (Rental) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("car_id"),
		index.Fields("deleted_at"),
		index.Fields("renter_id"),
		index.Fields("tenant_id"),
	}
}

// Annotations of the Rental.
func (Rental) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.RelayConnection(),
		entgql.QueryField(),
		entgql.Mutations(entgql.MutationCreate(), entgql.MutationUpdate()),
		// Add soft delete support
		entgql.MultiOrder(),
	}
}

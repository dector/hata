package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// House holds the schema definition for the House entity.
type House struct {
	ent.Schema
}

// Annotations of the House.
func (House) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "house"},
	}
}

// Fields of the House.
func (House) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			NotEmpty().
			Immutable().
			Comment("Sqids-generated house ID"),
		field.String("display_name").
			NotEmpty().
			Comment("Display name for the house"),
	}
}

// Edges of the House.
func (House) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("roles", HouseRole.Type),
		edge.To("devices", Device.Type),
		edge.To("shopping_lists", ShoppingList.Type),
	}
}

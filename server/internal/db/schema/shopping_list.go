package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// ShoppingList holds the schema definition for the ShoppingList entity.
type ShoppingList struct {
	ent.Schema
}

// Annotations of the ShoppingList.
func (ShoppingList) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "shopping__list"},
	}
}

// Fields of the ShoppingList.
func (ShoppingList) Fields() []ent.Field {
	return []ent.Field{
		field.String("house_id").
			NotEmpty().
			Comment("Foreign key to House"),
		field.String("uid").
			NotEmpty().
			Comment("External list UID (house-scoped)"),
		field.String("name").
			NotEmpty().
			Comment("Display name for the shopping list"),
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("When the shopping list was created"),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("When the shopping list was updated"),
	}
}

// Edges of the ShoppingList.
func (ShoppingList) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("house", House.Type).
			Ref("shopping_lists").
			Required().
			Unique().
			Field("house_id"),
		edge.To("items", ShoppingItem.Type),
	}
}

// Indexes of the ShoppingList.
func (ShoppingList) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("house_id", "uid").Unique(),
		index.Fields("house_id"),
	}
}

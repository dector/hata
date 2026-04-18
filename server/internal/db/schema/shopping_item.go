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

// ShoppingItem holds the schema definition for the ShoppingItem entity.
type ShoppingItem struct {
	ent.Schema
}

// Annotations of the ShoppingItem.
func (ShoppingItem) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "shopping__item"},
	}
}

// Fields of the ShoppingItem.
func (ShoppingItem) Fields() []ent.Field {
	return []ent.Field{
		field.Int("list_id").
			Comment("Foreign key to ShoppingList"),
		field.String("uid").
			NotEmpty().
			Comment("External item UID (list-scoped)"),
		field.String("name").
			NotEmpty().
			Comment("Item name"),
		field.Int("position").
			Comment("Stable server-side item ordering"),
		field.Time("checked_at").
			Optional().
			Nillable().
			Comment("When the item was checked"),
		field.Int("checked_by_user_id").
			Optional().
			Nillable().
			Comment("User who checked the item"),
		field.Time("deleted_at").
			Optional().
			Nillable().
			Comment("Soft-delete timestamp"),
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("When the shopping item was created"),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("When the shopping item was updated"),
	}
}

// Edges of the ShoppingItem.
func (ShoppingItem) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("list", ShoppingList.Type).
			Ref("items").
			Required().
			Unique().
			Field("list_id"),
		edge.To("checked_by_user", User.Type).
			Unique().
			Field("checked_by_user_id"),
	}
}

// Indexes of the ShoppingItem.
func (ShoppingItem) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("list_id", "uid").Unique(),
		index.Fields("list_id", "deleted_at", "position"),
		index.Fields("checked_by_user_id"),
	}
}

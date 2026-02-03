package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// HouseRole holds the schema definition for the HouseRole entity.
type HouseRole struct {
	ent.Schema
}

// Annotations of the HouseRole.
func (HouseRole) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "house__role"},
	}
}

// Fields of the HouseRole.
func (HouseRole) Fields() []ent.Field {
	return []ent.Field{
		field.Int("user_id").
			Comment("Foreign key to User"),
		field.String("house_id").
			Comment("Foreign key to House"),
		field.Enum("role").
			Values("admin", "owner", "habitant", "guest").
			Comment("Role for the user in the house"),
	}
}

// Edges of the HouseRole.
func (HouseRole) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("house_roles").
			Required().
			Unique().
			Field("user_id"),
		edge.From("house", House.Type).
			Ref("roles").
			Required().
			Unique().
			Field("house_id"),
	}
}

// Indexes of the HouseRole.
func (HouseRole) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("house_id", "user_id").Unique(),
		index.Fields("user_id"),
		index.Fields("house_id"),
	}
}

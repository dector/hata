package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// HouseMetaSetting holds plugin-like scoped configuration for a house.
type HouseMetaSetting struct {
	ent.Schema
}

// Annotations of the HouseMetaSetting.
func (HouseMetaSetting) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "house_ext_settings"},
	}
}

// Fields of the HouseMetaSetting.
func (HouseMetaSetting) Fields() []ent.Field {
	return []ent.Field{
		field.String("house_id").
			NotEmpty().
			Comment("Foreign key to House"),
		field.String("scope").
			NotEmpty().
			Comment("Settings scope, e.g. hata.ext.weather.v1"),
		field.String("key").
			NotEmpty().
			Comment("Setting key inside the scope"),
		field.String("value").
			Comment("Setting value stored as string"),
	}
}

// Edges of the HouseMetaSetting.
func (HouseMetaSetting) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("house", House.Type).
			Ref("meta_settings").
			Required().
			Unique().
			Field("house_id"),
	}
}

// Indexes of the HouseMetaSetting.
func (HouseMetaSetting) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("house_id", "scope", "key").Unique(),
		index.Fields("house_id"),
		index.Fields("scope"),
	}
}

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

// Session holds the schema definition for the Session entity.
type Session struct {
	ent.Schema
}

// Annotations of the Session.
func (Session) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "auth__Session"},
	}
}

// Fields of the Session.
func (Session) Fields() []ent.Field {
	return []ent.Field{
		field.Int("user_id").
			Comment("Foreign key to User"),
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("When the session was created"),
		field.Time("valid_until").
			Comment("Session expiration timestamp"),
		field.String("token").
			Unique().
			MaxLen(40).
			MinLen(40).
			Comment("40-character random session token"),
		field.Time("invalid_since").
			Optional().
			Nillable().
			Comment("If set, session is considered invalid"),
	}
}

// Edges of the Session.
func (Session) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user", User.Type).
			Required().
			Unique().
			Field("user_id"),
	}
}

// Indexes of the Session.
func (Session) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("token").Unique(),
		index.Fields("user_id"),
	}
}

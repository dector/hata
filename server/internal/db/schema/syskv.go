package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// SysKV holds the schema definition for the SysKV entity.
type SysKV struct {
	ent.Schema
}

// Annotations of the SysKV.
func (SysKV) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "sys__kv"},
	}
}

// Fields of the SysKV.
func (SysKV) Fields() []ent.Field {
	return []ent.Field{
		field.String("key"),
		field.String("value"),
	}
}

// Edges of the SysKV.
func (SysKV) Edges() []ent.Edge {
	return nil
}

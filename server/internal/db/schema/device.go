package schema

import (
	"fmt"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Device holds the schema definition for the Device entity.
type Device struct {
	ent.Schema
}

// Annotations of the Device.
func (Device) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "device"},
	}
}

// Fields of the Device.
func (Device) Fields() []ent.Field {
	return []ent.Field{
		field.String("device_id").
			NotEmpty().
			Immutable().
			MaxLen(64).
			Comment("Slug identifier for device"),
		field.String("name").
			NotEmpty().
			Comment("Display name for the device"),
		field.String("integration_id").
			NotEmpty().
			Comment("Integration identifier"),
		field.String("integration_data").
			Optional().
			Nillable().
			Comment("Integration-specific JSON data"),
			field.String("state").
			Default("").
			Validate(validateDeviceState).
			Comment("Current device state: on, off, or empty for unknown"),
		field.String("availability").
			Default("unknown").
			Validate(validateDeviceAvailability).
			Comment("Current device availability: online, offline, or unknown"),
		field.String("house_id").
			NotEmpty().
			Comment("Foreign key to House"),
	}
}

// Edges of the Device.
func (Device) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("house", House.Type).
			Ref("devices").
			Required().
			Unique().
			Field("house_id"),
	}
}

// Indexes of the Device.
func (Device) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("house_id", "device_id").Unique(),
		index.Fields("house_id"),
	}
}

func validateDeviceState(state string) error {
	switch state {
	case "", "on", "off":
		return nil
	default:
		return fmt.Errorf("invalid device state %q", state)
	}
}

func validateDeviceAvailability(availability string) error {
	switch availability {
	case "unknown", "online", "offline":
		return nil
	default:
		return fmt.Errorf("invalid device availability %q", availability)
	}
}

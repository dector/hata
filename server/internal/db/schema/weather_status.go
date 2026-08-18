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

// WeatherStatus holds the latest cached weather status for a house.
type WeatherStatus struct {
	ent.Schema
}

// Annotations of the WeatherStatus.
func (WeatherStatus) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "weather_status"},
	}
}

// Fields of the WeatherStatus.
func (WeatherStatus) Fields() []ent.Field {
	return []ent.Field{
		field.String("house_id").
			NotEmpty().
			Comment("Foreign key to House"),
		field.Float("temperature").
			Comment("Current temperature"),
		field.String("temperature_unit").
			NotEmpty().
			Comment("Temperature unit"),
		field.Int("condition_code").
			Comment("Provider weather condition code"),
		field.String("condition_text").
			NotEmpty().
			Comment("Human-readable weather condition"),
		field.String("condition_icon").
			NotEmpty().
			Comment("Internal weather condition icon key"),
		field.Float("humidity_percent").
			Optional().
			Nillable().
			Comment("Relative humidity percent"),
		field.Float("wind_speed").
			Optional().
			Nillable().
			Comment("Current wind speed"),
		field.String("wind_speed_unit").
			Optional().
			Nillable().
			Comment("Wind speed unit"),
		field.Time("observed_at").
			Comment("Provider observation timestamp"),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("Cache update timestamp"),
	}
}

// Edges of the WeatherStatus.
func (WeatherStatus) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("house", House.Type).
			Ref("weather_statuses").
			Required().
			Unique().
			Field("house_id"),
	}
}

// Indexes of the WeatherStatus.
func (WeatherStatus) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("house_id").Unique(),
	}
}

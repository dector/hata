package schema

import (
	"net"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// HouseDiscoveryNetwork holds CIDR ranges used for house device discovery.
type HouseDiscoveryNetwork struct {
	ent.Schema
}

// Annotations of the HouseDiscoveryNetwork.
func (HouseDiscoveryNetwork) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "house_discovery_network"},
	}
}

// Fields of the HouseDiscoveryNetwork.
func (HouseDiscoveryNetwork) Fields() []ent.Field {
	return []ent.Field{
		field.String("house_id").
			NotEmpty().
			Comment("Foreign key to House"),
		field.String("cidr").
			NotEmpty().
			Validate(validateDiscoveryCIDR).
			Comment("IPv4 CIDR scanned during discovery"),
		field.String("label").
			Default("").
			Comment("Human-readable discovery network label"),
	}
}

// Edges of the HouseDiscoveryNetwork.
func (HouseDiscoveryNetwork) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("house", House.Type).
			Ref("discovery_networks").
			Required().
			Unique().
			Field("house_id"),
	}
}

// Indexes of the HouseDiscoveryNetwork.
func (HouseDiscoveryNetwork) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("house_id", "cidr").Unique(),
		index.Fields("house_id"),
	}
}

func validateDiscoveryCIDR(cidr string) error {
	_, _, err := net.ParseCIDR(strings.TrimSpace(cidr))
	return err
}

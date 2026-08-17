package repo_test

import (
	"context"
	"testing"

	"hata/internal/db"
)

func TestHouseDiscoveryNetworkRepo_CreateListDelete(t *testing.T) {
	ctx := context.Background()
	dbInst, err := db.OpenTestDB(ctx)
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	defer dbInst.Close()

	repos := dbInst.Repos()
	if _, err := repos.House().Create(ctx, "H1", "Main Home"); err != nil {
		t.Fatalf("create house: %v", err)
	}

	network, err := repos.HouseDiscoveryNetwork().Create(ctx, "H1", "192.168.1.42/24", "Main WiFi")
	if err != nil {
		t.Fatalf("create network: %v", err)
	}
	if network.CIDR != "192.168.1.0/24" {
		t.Fatalf("expected normalized CIDR, got %q", network.CIDR)
	}

	networks, err := repos.HouseDiscoveryNetwork().ListByHouse(ctx, "H1")
	if err != nil {
		t.Fatalf("list networks: %v", err)
	}
	if len(networks) != 1 || networks[0].Label != "Main WiFi" {
		t.Fatalf("expected created network, got %#v", networks)
	}

	if err := repos.HouseDiscoveryNetwork().DeleteByID(ctx, "H1", network.ID); err != nil {
		t.Fatalf("delete network: %v", err)
	}
	networks, err = repos.HouseDiscoveryNetwork().ListByHouse(ctx, "H1")
	if err != nil {
		t.Fatalf("list after delete: %v", err)
	}
	if len(networks) != 0 {
		t.Fatalf("expected no networks after delete, got %#v", networks)
	}
}

func TestHouseDiscoveryNetworkRepo_CreateRejectsLargeNetworks(t *testing.T) {
	ctx := context.Background()
	dbInst, err := db.OpenTestDB(ctx)
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	defer dbInst.Close()

	repos := dbInst.Repos()
	if _, err := repos.House().Create(ctx, "H1", "Main Home"); err != nil {
		t.Fatalf("create house: %v", err)
	}
	if _, err := repos.HouseDiscoveryNetwork().Create(ctx, "H1", "192.168.0.0/16", "Too large"); err == nil {
		t.Fatalf("expected large CIDR to be rejected")
	}
}

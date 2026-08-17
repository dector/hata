package repo

import (
	"context"
	"fmt"
	"net"
	"strings"

	"hata/internal/orm"
	"hata/internal/orm/housediscoverynetwork"
)

// HouseDiscoveryNetworkRepo implements the HouseDiscoveryNetworkRepository interface.
type HouseDiscoveryNetworkRepo struct {
	client *orm.Client
}

// Create creates a discovery network for a house.
func (r *HouseDiscoveryNetworkRepo) Create(ctx context.Context, houseID, cidr, label string) (*HouseDiscoveryNetworkData, error) {
	cidr, err := normalizeCIDR(cidr)
	if err != nil {
		return nil, err
	}

	n, err := r.client.HouseDiscoveryNetwork.Create().
		SetHouseID(houseID).
		SetCidr(cidr).
		SetLabel(strings.TrimSpace(label)).
		Save(ctx)
	if err != nil {
		if orm.IsConstraintError(err) {
			return nil, fmt.Errorf("discovery network %q already exists in house %q", cidr, houseID)
		}
		return nil, fmt.Errorf("failed creating discovery network: %w", err)
	}
	return houseDiscoveryNetworkDataFromEnt(n), nil
}

// ListByHouse lists discovery networks for a house.
func (r *HouseDiscoveryNetworkRepo) ListByHouse(ctx context.Context, houseID string) ([]*HouseDiscoveryNetworkData, error) {
	networks, err := r.client.HouseDiscoveryNetwork.Query().
		Where(housediscoverynetwork.HouseIDEQ(houseID)).
		Order(housediscoverynetwork.ByCidr()).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed listing discovery networks for house %q: %w", houseID, err)
	}

	results := make([]*HouseDiscoveryNetworkData, len(networks))
	for i, n := range networks {
		results[i] = houseDiscoveryNetworkDataFromEnt(n)
	}
	return results, nil
}

// DeleteByID deletes a discovery network if it belongs to the house.
func (r *HouseDiscoveryNetworkRepo) DeleteByID(ctx context.Context, houseID string, id int) error {
	count, err := r.client.HouseDiscoveryNetwork.Delete().
		Where(housediscoverynetwork.IDEQ(id), housediscoverynetwork.HouseIDEQ(houseID)).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed deleting discovery network: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("discovery network %d not found in house %q", id, houseID)
	}
	return nil
}

func normalizeCIDR(cidr string) (string, error) {
	ip, ipNet, err := net.ParseCIDR(strings.TrimSpace(cidr))
	if err != nil {
		return "", fmt.Errorf("invalid CIDR %q", cidr)
	}
	ip4 := ip.To4()
	if ip4 == nil {
		return "", fmt.Errorf("invalid CIDR %q: only IPv4 networks are supported", cidr)
	}
	ones, bits := ipNet.Mask.Size()
	if bits != 32 || ones < 22 {
		return "", fmt.Errorf("invalid CIDR %q: discovery networks cannot be larger than /22", cidr)
	}
	ipNet.IP = ip4.Mask(ipNet.Mask)
	return ipNet.String(), nil
}

func houseDiscoveryNetworkDataFromEnt(n *orm.HouseDiscoveryNetwork) *HouseDiscoveryNetworkData {
	return &HouseDiscoveryNetworkData{
		ID:      n.ID,
		HouseID: n.HouseID,
		CIDR:    n.Cidr,
		Label:   n.Label,
	}
}

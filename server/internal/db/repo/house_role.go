package repo

import (
	"context"
	"fmt"

	"hata/internal/orm"
	"hata/internal/orm/houserole"
)

// HouseRoleRepo implements the HouseRoleRepository interface.
type HouseRoleRepo struct {
	client *orm.Client
}

// Assign assigns a role to a user for a house.
func (r *HouseRoleRepo) Assign(ctx context.Context, houseID string, userID int, role string) (*HouseRoleData, error) {
	hr, err := r.client.HouseRole.Create().
		SetHouseID(houseID).
		SetUserID(userID).
		SetRole(houserole.Role(role)).
		Save(ctx)

	if err != nil {
		if orm.IsConstraintError(err) {
			return nil, fmt.Errorf("user %d already assigned to house %q", userID, houseID)
		}
		return nil, fmt.Errorf("failed assigning house role: %w", err)
	}

	return &HouseRoleData{
		ID:      hr.ID,
		HouseID: hr.HouseID,
		UserID:  hr.UserID,
		Role:    string(hr.Role),
	}, nil
}

// ListByUser lists houses and roles for the given user.
func (r *HouseRoleRepo) ListByUser(ctx context.Context, userID int) ([]*HouseMembershipData, error) {
	roles, err := r.client.HouseRole.Query().
		Where(houserole.UserIDEQ(userID)).
		WithHouse().
		All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed listing house roles for user %d: %w", userID, err)
	}

	results := make([]*HouseMembershipData, 0, len(roles))
	for _, role := range roles {
		house := role.Edges.House
		if house == nil {
			return nil, fmt.Errorf("house role %d missing house edge", role.ID)
		}
		results = append(results, &HouseMembershipData{
			HouseID:     house.ID,
			DisplayName: house.DisplayName,
			Location:    house.Location,
			Role:        string(role.Role),
		})
	}

	return results, nil
}

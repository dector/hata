package repo

import (
	"context"
	"fmt"

	"hata/internal/orm"
	"hata/internal/orm/house"
)

// HouseRepo implements the HouseRepository interface.
type HouseRepo struct {
	client *orm.Client
}

// Create creates a new house with the given ID and display name.
func (r *HouseRepo) Create(ctx context.Context, id, displayName string) (*HouseData, error) {
	h, err := r.client.House.Create().
		SetID(id).
		SetDisplayName(displayName).
		Save(ctx)

	if err != nil {
		if orm.IsConstraintError(err) {
			return nil, fmt.Errorf("house %q already exists", id)
		}
		return nil, fmt.Errorf("failed creating house: %w", err)
	}

	return &HouseData{
		ID:          h.ID,
		DisplayName: h.DisplayName,
	}, nil
}

// GetByID retrieves a house by its ID.
func (r *HouseRepo) GetByID(ctx context.Context, id string) (*HouseData, error) {
	h, err := r.client.House.Query().
		Where(house.IDEQ(id)).
		Only(ctx)

	if err != nil {
		if orm.IsNotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed querying house by id %q: %w", id, err)
	}

	return &HouseData{
		ID:          h.ID,
		DisplayName: h.DisplayName,
	}, nil
}

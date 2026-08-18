package repo

import (
	"context"
	"errors"
	"fmt"

	"hata/internal/orm"
	"hata/internal/orm/house"
)

// ErrHouseNotFound is returned when a house does not exist.
var ErrHouseNotFound = errors.New("house not found")

// HouseRepo implements the HouseRepository interface.
type HouseRepo struct {
	client *orm.Client
}

// Create creates a new house with the given ID and display name.
func (r *HouseRepo) Create(ctx context.Context, id, displayName string) (*HouseData, error) {
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed starting transaction for house creation: %w", err)
	}

	h, err := tx.House.Create().
		SetID(id).
		SetDisplayName(displayName).
		Save(ctx)
	if err != nil {
		_ = tx.Rollback()
		if orm.IsConstraintError(err) {
			return nil, fmt.Errorf("house %q already exists", id)
		}
		return nil, fmt.Errorf("failed creating house: %w", err)
	}

	_, err = tx.ShoppingList.Create().
		SetHouseID(h.ID).
		SetUID(DefaultShoppingListUID).
		SetName(DefaultShoppingListName).
		Save(ctx)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("failed creating default shopping list for house %q: %w", id, err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed committing house creation transaction: %w", err)
	}

	return houseDataFromEnt(h), nil
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

	return houseDataFromEnt(h), nil
}

// ListAll lists all houses.
func (r *HouseRepo) ListAll(ctx context.Context) ([]*HouseData, error) {
	houses, err := r.client.House.Query().
		Order(orm.Asc(house.FieldID)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed listing houses: %w", err)
	}
	results := make([]*HouseData, 0, len(houses))
	for _, h := range houses {
		results = append(results, houseDataFromEnt(h))
	}
	return results, nil
}

// UpdateLocation updates a house location.
func (r *HouseRepo) UpdateLocation(ctx context.Context, id string, location *string) error {
	update := r.client.House.Update().
		Where(house.IDEQ(id))
	if location == nil {
		update.ClearLocation()
	} else {
		update.SetLocation(*location)
	}

	count, err := update.Save(ctx)
	if err != nil {
		return fmt.Errorf("failed updating house %q location: %w", id, err)
	}
	if count == 0 {
		return fmt.Errorf("%w: %q", ErrHouseNotFound, id)
	}
	return nil
}

func houseDataFromEnt(h *orm.House) *HouseData {
	if h == nil {
		return nil
	}
	return &HouseData{
		ID:          h.ID,
		DisplayName: h.DisplayName,
		Location:    h.Location,
	}
}

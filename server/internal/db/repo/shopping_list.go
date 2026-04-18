package repo

import (
	"context"
	"errors"
	"fmt"

	"hata/internal/orm"
	"hata/internal/orm/shoppinglist"
)

const (
	DefaultShoppingListUID  = "default"
	DefaultShoppingListName = "default"
)

var ErrShoppingListNotFound = errors.New("shopping list not found")

// ShoppingListRepo implements the ShoppingListRepository interface.
type ShoppingListRepo struct {
	client *orm.Client
}

// Create creates a shopping list in a house.
func (r *ShoppingListRepo) Create(ctx context.Context, houseID, uid, name string) (*ShoppingListData, error) {
	list, err := r.client.ShoppingList.Create().
		SetHouseID(houseID).
		SetUID(uid).
		SetName(name).
		Save(ctx)
	if err != nil {
		if orm.IsConstraintError(err) {
			return nil, fmt.Errorf("shopping list %q already exists in house %q", uid, houseID)
		}
		return nil, fmt.Errorf("failed creating shopping list: %w", err)
	}

	return shoppingListDataFromEnt(list), nil
}

// ListByHouse lists shopping lists for a house.
func (r *ShoppingListRepo) ListByHouse(ctx context.Context, houseID string) ([]*ShoppingListData, error) {
	lists, err := r.client.ShoppingList.Query().
		Where(shoppinglist.HouseIDEQ(houseID)).
		Order(shoppinglist.ByCreatedAt(), shoppinglist.ByUID()).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed listing shopping lists for house %q: %w", houseID, err)
	}

	result := make([]*ShoppingListData, len(lists))
	for i, list := range lists {
		result[i] = shoppingListDataFromEnt(list)
	}

	return result, nil
}

// GetByHouseAndUID retrieves a shopping list by house and list UID.
func (r *ShoppingListRepo) GetByHouseAndUID(ctx context.Context, houseID, uid string) (*ShoppingListData, error) {
	list, err := r.client.ShoppingList.Query().
		Where(shoppinglist.HouseIDEQ(houseID), shoppinglist.UIDEQ(uid)).
		Only(ctx)
	if err != nil {
		if orm.IsNotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed querying shopping list %q for house %q: %w", uid, houseID, err)
	}

	return shoppingListDataFromEnt(list), nil
}

func shoppingListDataFromEnt(list *orm.ShoppingList) *ShoppingListData {
	return &ShoppingListData{
		ID:        list.ID,
		HouseID:   list.HouseID,
		UID:       list.UID,
		Name:      list.Name,
		CreatedAt: list.CreatedAt,
		UpdatedAt: list.UpdatedAt,
	}
}

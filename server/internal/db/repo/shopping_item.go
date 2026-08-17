package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	entsql "entgo.io/ent/dialect/sql"

	"hata/internal/orm"
	"hata/internal/orm/predicate"
	"hata/internal/orm/shoppingitem"
	"hata/internal/orm/shoppinglist"
)

var ErrShoppingItemNotFound = errors.New("shopping item not found")

// ShoppingItemRepo implements the ShoppingItemRepository interface.
type ShoppingItemRepo struct {
	client *orm.Client
}

// Create creates a shopping item in the specified list and appends it to the end.
func (r *ShoppingItemRepo) Create(ctx context.Context, houseID, listUID, itemUID, name string) (*ShoppingItemData, error) {
	list, err := r.client.ShoppingList.Query().
		Where(shoppinglist.HouseIDEQ(houseID), shoppinglist.UIDEQ(listUID)).
		Only(ctx)
	if err != nil {
		if orm.IsNotFound(err) {
			return nil, fmt.Errorf("%w: %q in house %q", ErrShoppingListNotFound, listUID, houseID)
		}
		return nil, fmt.Errorf("failed loading shopping list %q in house %q: %w", listUID, houseID, err)
	}

	nextPosition, err := r.nextPositionInList(ctx, list.ID)
	if err != nil {
		return nil, err
	}

	item, err := r.client.ShoppingItem.Create().
		SetListID(list.ID).
		SetUID(itemUID).
		SetName(name).
		SetPosition(nextPosition).
		Save(ctx)
	if err != nil {
		if orm.IsConstraintError(err) {
			return nil, fmt.Errorf("shopping item %q already exists in list %q", itemUID, listUID)
		}
		return nil, fmt.Errorf("failed creating shopping item: %w", err)
	}

	return shoppingItemDataFromEnt(item), nil
}

// ListByHouseAndListUID lists shopping items for a list.
func (r *ShoppingItemRepo) ListByHouseAndListUID(ctx context.Context, houseID, listUID string, includeDeleted bool) ([]*ShoppingItemData, error) {
	predicates := []predicate.ShoppingItem{
		shoppingitem.HasListWith(
			shoppinglist.HouseIDEQ(houseID),
			shoppinglist.UIDEQ(listUID),
		),
	}
	if !includeDeleted {
		predicates = append(predicates, shoppingitem.DeletedAtIsNil())
	}

	items, err := r.client.ShoppingItem.Query().
		Where(predicates...).
		Order(shoppingitem.ByPosition(), shoppingitem.ByCreatedAt()).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed listing shopping items for list %q in house %q: %w", listUID, houseID, err)
	}

	result := make([]*ShoppingItemData, len(items))
	for i, item := range items {
		result[i] = shoppingItemDataFromEnt(item)
	}

	return result, nil
}

// GetByHouseAndListAndUID retrieves a shopping item by house/list/item UIDs.
func (r *ShoppingItemRepo) GetByHouseAndListAndUID(ctx context.Context, houseID, listUID, itemUID string) (*ShoppingItemData, error) {
	item, err := r.client.ShoppingItem.Query().
		Where(
			shoppingitem.UIDEQ(itemUID),
			shoppingitem.HasListWith(
				shoppinglist.HouseIDEQ(houseID),
				shoppinglist.UIDEQ(listUID),
			),
		).
		Only(ctx)
	if err != nil {
		if orm.IsNotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed loading shopping item %q in list %q (house %q): %w", itemUID, listUID, houseID, err)
	}

	return shoppingItemDataFromEnt(item), nil
}

// UpdateName updates shopping item name.
func (r *ShoppingItemRepo) UpdateName(ctx context.Context, houseID, listUID, itemUID, name string) error {
	count, err := r.client.ShoppingItem.Update().
		Where(
			shoppingitem.UIDEQ(itemUID),
			shoppingitem.HasListWith(
				shoppinglist.HouseIDEQ(houseID),
				shoppinglist.UIDEQ(listUID),
			),
		).
		SetName(name).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed updating shopping item name: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("%w: %q in list %q (house %q)", ErrShoppingItemNotFound, itemUID, listUID, houseID)
	}

	return nil
}

// SetChecked sets checked state fields.
func (r *ShoppingItemRepo) SetChecked(ctx context.Context, houseID, listUID, itemUID string, checkedAt *time.Time, checkedByUserID *int) error {
	var moveToPosition *int
	if checkedAt != nil {
		item, err := r.client.ShoppingItem.Query().
			Where(
				shoppingitem.UIDEQ(itemUID),
				shoppingitem.HasListWith(
					shoppinglist.HouseIDEQ(houseID),
					shoppinglist.UIDEQ(listUID),
				),
			).
			First(ctx)
		if err != nil {
			if orm.IsNotFound(err) {
				return fmt.Errorf("%w: %q in list %q (house %q)", ErrShoppingItemNotFound, itemUID, listUID, houseID)
			}
			return fmt.Errorf("failed loading shopping item before check state update: %w", err)
		}
		nextPosition, err := r.nextPositionInList(ctx, item.ListID)
		if err != nil {
			return err
		}
		moveToPosition = &nextPosition
	}

	update := r.client.ShoppingItem.Update().
		Where(
			shoppingitem.UIDEQ(itemUID),
			shoppingitem.HasListWith(
				shoppinglist.HouseIDEQ(houseID),
				shoppinglist.UIDEQ(listUID),
			),
		)

	if checkedAt == nil {
		update = update.ClearCheckedAt()
	} else {
		update = update.SetCheckedAt(*checkedAt)
	}

	if checkedByUserID == nil {
		update = update.ClearCheckedByUserID()
	} else {
		update = update.SetCheckedByUserID(*checkedByUserID)
	}
	if moveToPosition != nil {
		update = update.SetPosition(*moveToPosition)
	}

	count, err := update.Save(ctx)
	if err != nil {
		return fmt.Errorf("failed updating shopping item check state: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("%w: %q in list %q (house %q)", ErrShoppingItemNotFound, itemUID, listUID, houseID)
	}

	return nil
}

// SoftDelete marks shopping item as deleted.
func (r *ShoppingItemRepo) SoftDelete(ctx context.Context, houseID, listUID, itemUID string, deletedAt time.Time) error {
	count, err := r.client.ShoppingItem.Update().
		Where(
			shoppingitem.UIDEQ(itemUID),
			shoppingitem.HasListWith(
				shoppinglist.HouseIDEQ(houseID),
				shoppinglist.UIDEQ(listUID),
			),
		).
		SetDeletedAt(deletedAt).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed soft-deleting shopping item: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("%w: %q in list %q (house %q)", ErrShoppingItemNotFound, itemUID, listUID, houseID)
	}

	return nil
}

func (r *ShoppingItemRepo) nextPositionInList(ctx context.Context, listID int) (int, error) {
	last, err := r.client.ShoppingItem.Query().
		Where(shoppingitem.ListIDEQ(listID)).
		Order(shoppingitem.ByPosition(entsql.OrderDesc()), shoppingitem.ByID(entsql.OrderDesc())).
		First(ctx)
	if err != nil {
		if orm.IsNotFound(err) {
			return 0, nil
		}
		return 0, fmt.Errorf("failed computing next shopping item position for list %d: %w", listID, err)
	}

	return last.Position + 1, nil
}

func shoppingItemDataFromEnt(item *orm.ShoppingItem) *ShoppingItemData {
	return &ShoppingItemData{
		ID:              item.ID,
		ListID:          item.ListID,
		UID:             item.UID,
		Name:            item.Name,
		Position:        item.Position,
		CheckedAt:       item.CheckedAt,
		CheckedByUserID: item.CheckedByUserID,
		DeletedAt:       item.DeletedAt,
		CreatedAt:       item.CreatedAt,
		UpdatedAt:       item.UpdatedAt,
	}
}

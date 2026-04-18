package repo_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"hata/internal/db"
	repo "hata/internal/db/repo"
)

func TestHouseCreate_CreatesDefaultShoppingList(t *testing.T) {
	repos := openTestRepos(t)
	ctx := context.Background()

	house, err := repos.House().Create(ctx, "house-1", "Main House")
	if err != nil {
		t.Fatalf("create house: %v", err)
	}

	lists, err := repos.ShoppingList().ListByHouse(ctx, house.ID)
	if err != nil {
		t.Fatalf("list shopping lists: %v", err)
	}

	if len(lists) != 1 {
		t.Fatalf("expected 1 shopping list, got %d", len(lists))
	}

	if lists[0].UID != repo.DefaultShoppingListUID {
		t.Fatalf("expected default list uid %q, got %q", repo.DefaultShoppingListUID, lists[0].UID)
	}
	if lists[0].Name != repo.DefaultShoppingListName {
		t.Fatalf("expected default list name %q, got %q", repo.DefaultShoppingListName, lists[0].Name)
	}
}

func TestShoppingListRepo_CreateAndGetByHouseAndUID(t *testing.T) {
	repos := openTestRepos(t)
	ctx := context.Background()

	house, err := repos.House().Create(ctx, "house-1", "Main House")
	if err != nil {
		t.Fatalf("create house: %v", err)
	}

	created, err := repos.ShoppingList().Create(ctx, house.ID, "weekly", "Weekly")
	if err != nil {
		t.Fatalf("create shopping list: %v", err)
	}
	if created.UID != "weekly" {
		t.Fatalf("expected created uid weekly, got %q", created.UID)
	}

	got, err := repos.ShoppingList().GetByHouseAndUID(ctx, house.ID, "weekly")
	if err != nil {
		t.Fatalf("get shopping list: %v", err)
	}
	if got == nil {
		t.Fatal("expected shopping list, got nil")
	}
	if got.Name != "Weekly" {
		t.Fatalf("expected name Weekly, got %q", got.Name)
	}

	_, err = repos.ShoppingList().Create(ctx, house.ID, "weekly", "Weekly duplicate")
	if err == nil {
		t.Fatal("expected duplicate list create to fail")
	}
}

func TestShoppingItemRepo_CreateSoftDeleteAndFilter(t *testing.T) {
	repos := openTestRepos(t)
	ctx := context.Background()

	house, err := repos.House().Create(ctx, "house-1", "Main House")
	if err != nil {
		t.Fatalf("create house: %v", err)
	}

	item1, err := repos.ShoppingItem().Create(ctx, house.ID, repo.DefaultShoppingListUID, "item-1", "Milk")
	if err != nil {
		t.Fatalf("create item1: %v", err)
	}
	item2, err := repos.ShoppingItem().Create(ctx, house.ID, repo.DefaultShoppingListUID, "item-2", "Eggs")
	if err != nil {
		t.Fatalf("create item2: %v", err)
	}

	if item1.Position != 0 || item2.Position != 1 {
		t.Fatalf("expected positions 0 and 1, got %d and %d", item1.Position, item2.Position)
	}

	if err := repos.ShoppingItem().SoftDelete(ctx, house.ID, repo.DefaultShoppingListUID, "item-1", time.Now().UTC()); err != nil {
		t.Fatalf("soft delete item1: %v", err)
	}

	active, err := repos.ShoppingItem().ListByHouseAndListUID(ctx, house.ID, repo.DefaultShoppingListUID, false)
	if err != nil {
		t.Fatalf("list active items: %v", err)
	}
	if len(active) != 1 || active[0].UID != "item-2" {
		t.Fatalf("expected only item-2 active, got %+v", active)
	}

	all, err := repos.ShoppingItem().ListByHouseAndListUID(ctx, house.ID, repo.DefaultShoppingListUID, true)
	if err != nil {
		t.Fatalf("list all items: %v", err)
	}
	if len(all) != 2 {
		t.Fatalf("expected 2 items including deleted, got %d", len(all))
	}
}

func TestShoppingItemRepo_SetCheckedAndUncheck(t *testing.T) {
	repos := openTestRepos(t)
	ctx := context.Background()

	house, err := repos.House().Create(ctx, "house-1", "Main House")
	if err != nil {
		t.Fatalf("create house: %v", err)
	}
	user, err := repos.User().Create(ctx, "alice@example.com", "hash", "Alice")
	if err != nil {
		t.Fatalf("create user: %v", err)
	}
	_, err = repos.ShoppingItem().Create(ctx, house.ID, repo.DefaultShoppingListUID, "item-1", "Milk")
	if err != nil {
		t.Fatalf("create item: %v", err)
	}

	checkedAt := time.Now().UTC().Truncate(time.Second)
	checkedBy := user.ID
	if err := repos.ShoppingItem().SetChecked(ctx, house.ID, repo.DefaultShoppingListUID, "item-1", &checkedAt, &checkedBy); err != nil {
		t.Fatalf("set checked: %v", err)
	}

	item, err := repos.ShoppingItem().GetByHouseAndListAndUID(ctx, house.ID, repo.DefaultShoppingListUID, "item-1")
	if err != nil {
		t.Fatalf("get item: %v", err)
	}
	if item.CheckedAt == nil || !item.CheckedAt.Equal(checkedAt) {
		t.Fatalf("expected checkedAt %v, got %v", checkedAt, item.CheckedAt)
	}
	if item.CheckedByUserID == nil || *item.CheckedByUserID != user.ID {
		t.Fatalf("expected checkedByUserID %d, got %v", user.ID, item.CheckedByUserID)
	}

	if err := repos.ShoppingItem().SetChecked(ctx, house.ID, repo.DefaultShoppingListUID, "item-1", nil, nil); err != nil {
		t.Fatalf("uncheck: %v", err)
	}

	item, err = repos.ShoppingItem().GetByHouseAndListAndUID(ctx, house.ID, repo.DefaultShoppingListUID, "item-1")
	if err != nil {
		t.Fatalf("get item after uncheck: %v", err)
	}
	if item.CheckedAt != nil || item.CheckedByUserID != nil {
		t.Fatalf("expected item to be unchecked, got checkedAt=%v checkedByUserID=%v", item.CheckedAt, item.CheckedByUserID)
	}
}

func TestShoppingItemRepo_UpdateMissing_ReturnsNotFound(t *testing.T) {
	repos := openTestRepos(t)
	ctx := context.Background()

	house, err := repos.House().Create(ctx, "house-1", "Main House")
	if err != nil {
		t.Fatalf("create house: %v", err)
	}

	err = repos.ShoppingItem().UpdateName(ctx, house.ID, repo.DefaultShoppingListUID, "missing", "New name")
	if !errors.Is(err, repo.ErrShoppingItemNotFound) {
		t.Fatalf("expected ErrShoppingItemNotFound, got %v", err)
	}
}

func openTestRepos(t *testing.T) db.Repositories {
	t.Helper()

	dbInst, err := db.OpenTestDB(context.Background())
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	t.Cleanup(func() {
		if err := dbInst.Close(); err != nil {
			t.Fatalf("close test db: %v", err)
		}
	})

	return dbInst.Repos()
}

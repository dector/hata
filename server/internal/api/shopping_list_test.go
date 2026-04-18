package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	repo "hata/internal/db/repo"

	"github.com/go-chi/chi/v5"
)

func TestShoppingList_ListByHouse_Unauthorized(t *testing.T) {
	repos, cleanup := setupAuthTest(t)
	defer cleanup()

	handler := NewShoppingListHandler(repos)
	router := chi.NewRouter()
	router.Get("/api/latest/house/{houseId}/shopping-list", handler.ListByHouse)

	req := httptest.NewRequest(http.MethodGet, "/api/latest/house/house-1/shopping-list", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("Expected status 401, got %d", w.Code)
	}
}

func TestShoppingList_ListByHouse_GuestForbidden(t *testing.T) {
	repos, cleanup := setupAuthTest(t)
	defer cleanup()

	ctx := context.Background()
	user := createTestUser(t, repos, "guest@example.com")
	token := createTestSession(t, repos, user.ID, strings.Repeat("s", 40))

	house, err := repos.House().Create(ctx, "house-1", "Main House")
	if err != nil {
		t.Fatalf("Failed to create house: %v", err)
	}
	if _, err := repos.HouseRole().Assign(ctx, house.ID, user.ID, "guest"); err != nil {
		t.Fatalf("Failed to assign role: %v", err)
	}

	handler := NewShoppingListHandler(repos)
	router := chi.NewRouter()
	router.Get("/api/latest/house/{houseId}/shopping-list", handler.ListByHouse)

	req := httptest.NewRequest(http.MethodGet, "/api/latest/house/house-1/shopping-list", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("Expected status 403, got %d", w.Code)
	}
}

func TestShoppingList_ListByHouse_OK(t *testing.T) {
	repos, cleanup := setupAuthTest(t)
	defer cleanup()

	ctx := context.Background()
	user := createTestUser(t, repos, "owner@example.com")
	token := createTestSession(t, repos, user.ID, strings.Repeat("t", 40))

	house, err := repos.House().Create(ctx, "house-1", "Main House")
	if err != nil {
		t.Fatalf("Failed to create house: %v", err)
	}
	if _, err := repos.HouseRole().Assign(ctx, house.ID, user.ID, "owner"); err != nil {
		t.Fatalf("Failed to assign role: %v", err)
	}

	handler := NewShoppingListHandler(repos)
	router := chi.NewRouter()
	router.Get("/api/latest/house/{houseId}/shopping-list", handler.ListByHouse)

	req := httptest.NewRequest(http.MethodGet, "/api/latest/house/house-1/shopping-list", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}

	var resp ShoppingListListResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(resp.ShoppingLists) != 1 {
		t.Fatalf("Expected 1 shopping list, got %d", len(resp.ShoppingLists))
	}
	if resp.ShoppingLists[0].UID != repo.DefaultShoppingListUID {
		t.Fatalf("Expected default uid %q, got %q", repo.DefaultShoppingListUID, resp.ShoppingLists[0].UID)
	}
}

func TestShoppingItem_CreateDeleteAndIncludeDeleted(t *testing.T) {
	repos, cleanup := setupAuthTest(t)
	defer cleanup()

	ctx := context.Background()
	user := createTestUser(t, repos, "habit@example.com")
	token := createTestSession(t, repos, user.ID, strings.Repeat("u", 40))

	house, err := repos.House().Create(ctx, "house-1", "Main House")
	if err != nil {
		t.Fatalf("Failed to create house: %v", err)
	}
	if _, err := repos.HouseRole().Assign(ctx, house.ID, user.ID, "habitant"); err != nil {
		t.Fatalf("Failed to assign role: %v", err)
	}

	handler := NewShoppingListHandler(repos)
	router := chi.NewRouter()
	router.Post("/api/latest/house/{houseId}/shopping-list/{listId}/item", handler.CreateItem)
	router.Get("/api/latest/house/{houseId}/shopping-list/{listId}/item", handler.ListItems)
	router.Delete("/api/latest/house/{houseId}/shopping-list/{listId}/item/{itemId}", handler.DeleteItem)

	// Create item
	createReq := httptest.NewRequest(http.MethodPost, "/api/latest/house/house-1/shopping-list/default/item", strings.NewReader(`{"name":"  Milk  "}`))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("Authorization", "Bearer "+token)
	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createReq)

	if createW.Code != http.StatusCreated {
		t.Fatalf("Expected status 201, got %d", createW.Code)
	}

	var created ShoppingItemResponse
	if err := json.NewDecoder(createW.Body).Decode(&created); err != nil {
		t.Fatalf("Failed to decode create response: %v", err)
	}
	if created.Item.UID == "" {
		t.Fatal("Expected generated item uid")
	}
	if created.Item.Name != "Milk" {
		t.Fatalf("Expected trimmed item name Milk, got %q", created.Item.Name)
	}

	// List non-deleted: should contain 1
	listReq := httptest.NewRequest(http.MethodGet, "/api/latest/house/house-1/shopping-list/default/item", nil)
	listReq.Header.Set("Authorization", "Bearer "+token)
	listW := httptest.NewRecorder()
	router.ServeHTTP(listW, listReq)
	if listW.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", listW.Code)
	}
	var listed ShoppingItemListResponse
	if err := json.NewDecoder(listW.Body).Decode(&listed); err != nil {
		t.Fatalf("Failed to decode list response: %v", err)
	}
	if len(listed.Items) != 1 {
		t.Fatalf("Expected 1 item, got %d", len(listed.Items))
	}

	// Delete item (soft delete)
	deleteReq := httptest.NewRequest(http.MethodDelete, "/api/latest/house/house-1/shopping-list/default/item/"+created.Item.UID, nil)
	deleteReq.Header.Set("Authorization", "Bearer "+token)
	deleteW := httptest.NewRecorder()
	router.ServeHTTP(deleteW, deleteReq)
	if deleteW.Code != http.StatusNoContent {
		t.Fatalf("Expected status 204, got %d", deleteW.Code)
	}

	// List non-deleted: should be empty
	listAfterReq := httptest.NewRequest(http.MethodGet, "/api/latest/house/house-1/shopping-list/default/item", nil)
	listAfterReq.Header.Set("Authorization", "Bearer "+token)
	listAfterW := httptest.NewRecorder()
	router.ServeHTTP(listAfterW, listAfterReq)
	if listAfterW.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", listAfterW.Code)
	}
	var listedAfter ShoppingItemListResponse
	if err := json.NewDecoder(listAfterW.Body).Decode(&listedAfter); err != nil {
		t.Fatalf("Failed to decode list-after response: %v", err)
	}
	if len(listedAfter.Items) != 0 {
		t.Fatalf("Expected 0 active items, got %d", len(listedAfter.Items))
	}

	// List with include_deleted=true: should contain deleted item
	listDeletedReq := httptest.NewRequest(http.MethodGet, "/api/latest/house/house-1/shopping-list/default/item?include_deleted=true", nil)
	listDeletedReq.Header.Set("Authorization", "Bearer "+token)
	listDeletedW := httptest.NewRecorder()
	router.ServeHTTP(listDeletedW, listDeletedReq)
	if listDeletedW.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", listDeletedW.Code)
	}
	var withDeleted ShoppingItemListResponse
	if err := json.NewDecoder(listDeletedW.Body).Decode(&withDeleted); err != nil {
		t.Fatalf("Failed to decode include-deleted response: %v", err)
	}
	if len(withDeleted.Items) != 1 {
		t.Fatalf("Expected 1 item when include_deleted=true, got %d", len(withDeleted.Items))
	}
	if withDeleted.Items[0].DeletedAt == nil {
		t.Fatal("Expected deletedAt to be set for soft-deleted item")
	}
}

func TestShoppingItem_CheckAndUncheck(t *testing.T) {
	repos, cleanup := setupAuthTest(t)
	defer cleanup()

	ctx := context.Background()
	user := createTestUser(t, repos, "checker@example.com")
	token := createTestSession(t, repos, user.ID, strings.Repeat("v", 40))

	house, err := repos.House().Create(ctx, "house-1", "Main House")
	if err != nil {
		t.Fatalf("Failed to create house: %v", err)
	}
	if _, err := repos.HouseRole().Assign(ctx, house.ID, user.ID, "owner"); err != nil {
		t.Fatalf("Failed to assign role: %v", err)
	}
	created, err := repos.ShoppingItem().Create(ctx, house.ID, repo.DefaultShoppingListUID, "item-1", "Milk")
	if err != nil {
		t.Fatalf("Failed to seed shopping item: %v", err)
	}

	handler := NewShoppingListHandler(repos)
	router := chi.NewRouter()
	router.Patch("/api/latest/house/{houseId}/shopping-list/{listId}/item/{itemId}/check", handler.SetItemChecked)

	checkReq := httptest.NewRequest(http.MethodPatch, "/api/latest/house/house-1/shopping-list/default/item/"+created.UID+"/check", strings.NewReader(`{"checked":true}`))
	checkReq.Header.Set("Content-Type", "application/json")
	checkReq.Header.Set("Authorization", "Bearer "+token)
	checkW := httptest.NewRecorder()
	router.ServeHTTP(checkW, checkReq)
	if checkW.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", checkW.Code)
	}
	var checked ShoppingItemResponse
	if err := json.NewDecoder(checkW.Body).Decode(&checked); err != nil {
		t.Fatalf("Failed to decode checked response: %v", err)
	}
	if checked.Item.CheckedAt == nil {
		t.Fatal("Expected checkedAt to be set")
	}
	if checked.Item.CheckedByUserID == nil || *checked.Item.CheckedByUserID != user.ID {
		t.Fatalf("Expected checkedByUserId=%d, got %v", user.ID, checked.Item.CheckedByUserID)
	}

	uncheckReq := httptest.NewRequest(http.MethodPatch, "/api/latest/house/house-1/shopping-list/default/item/"+created.UID+"/check", strings.NewReader(`{"checked":false}`))
	uncheckReq.Header.Set("Content-Type", "application/json")
	uncheckReq.Header.Set("Authorization", "Bearer "+token)
	uncheckW := httptest.NewRecorder()
	router.ServeHTTP(uncheckW, uncheckReq)
	if uncheckW.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", uncheckW.Code)
	}
	var unchecked ShoppingItemResponse
	if err := json.NewDecoder(uncheckW.Body).Decode(&unchecked); err != nil {
		t.Fatalf("Failed to decode unchecked response: %v", err)
	}
	if unchecked.Item.CheckedAt != nil || unchecked.Item.CheckedByUserID != nil {
		t.Fatalf("Expected item to be unchecked, got checkedAt=%v checkedByUserID=%v", unchecked.Item.CheckedAt, unchecked.Item.CheckedByUserID)
	}
}

func TestShoppingItem_MutateSoftDeleted_ReturnsNotFound(t *testing.T) {
	repos, cleanup := setupAuthTest(t)
	defer cleanup()

	ctx := context.Background()
	user := createTestUser(t, repos, "mutator@example.com")
	token := createTestSession(t, repos, user.ID, strings.Repeat("w", 40))

	house, err := repos.House().Create(ctx, "house-1", "Main House")
	if err != nil {
		t.Fatalf("Failed to create house: %v", err)
	}
	if _, err := repos.HouseRole().Assign(ctx, house.ID, user.ID, "owner"); err != nil {
		t.Fatalf("Failed to assign role: %v", err)
	}
	item, err := repos.ShoppingItem().Create(ctx, house.ID, repo.DefaultShoppingListUID, "item-1", "Milk")
	if err != nil {
		t.Fatalf("Failed to seed shopping item: %v", err)
	}
	if err := repos.ShoppingItem().SoftDelete(ctx, house.ID, repo.DefaultShoppingListUID, item.UID, item.CreatedAt); err != nil {
		t.Fatalf("Failed to soft-delete seeded item: %v", err)
	}

	handler := NewShoppingListHandler(repos)
	router := chi.NewRouter()
	router.Patch("/api/latest/house/{houseId}/shopping-list/{listId}/item/{itemId}", handler.UpdateItem)
	router.Patch("/api/latest/house/{houseId}/shopping-list/{listId}/item/{itemId}/check", handler.SetItemChecked)
	router.Delete("/api/latest/house/{houseId}/shopping-list/{listId}/item/{itemId}", handler.DeleteItem)

	t.Run("patch name", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/api/latest/house/house-1/shopping-list/default/item/"+item.UID, strings.NewReader(`{"name":"Bread"}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusNotFound {
			t.Fatalf("Expected status 404, got %d", w.Code)
		}
	})

	t.Run("patch check", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/api/latest/house/house-1/shopping-list/default/item/"+item.UID+"/check", strings.NewReader(`{"checked":true}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusNotFound {
			t.Fatalf("Expected status 404, got %d", w.Code)
		}
	})

	t.Run("delete", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/latest/house/house-1/shopping-list/default/item/"+item.UID, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusNotFound {
			t.Fatalf("Expected status 404, got %d", w.Code)
		}
	})
}

func TestShoppingItem_List_InvalidIncludeDeletedParam(t *testing.T) {
	repos, cleanup := setupAuthTest(t)
	defer cleanup()

	ctx := context.Background()
	user := createTestUser(t, repos, "reader@example.com")
	token := createTestSession(t, repos, user.ID, strings.Repeat("x", 40))

	house, err := repos.House().Create(ctx, "house-1", "Main House")
	if err != nil {
		t.Fatalf("Failed to create house: %v", err)
	}
	if _, err := repos.HouseRole().Assign(ctx, house.ID, user.ID, "owner"); err != nil {
		t.Fatalf("Failed to assign role: %v", err)
	}

	handler := NewShoppingListHandler(repos)
	router := chi.NewRouter()
	router.Get("/api/latest/house/{houseId}/shopping-list/{listId}/item", handler.ListItems)

	req := httptest.NewRequest(http.MethodGet, "/api/latest/house/house-1/shopping-list/default/item?include_deleted=not-bool", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400, got %d", w.Code)
	}
}

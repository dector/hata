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

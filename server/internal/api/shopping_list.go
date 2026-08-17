package api

import "hata/internal/db"

// ShoppingListHandler provides shopping list and item endpoints.
type ShoppingListHandler struct {
	repos db.Repositories
}

// NewShoppingListHandler creates a new ShoppingListHandler.
func NewShoppingListHandler(repos db.Repositories) *ShoppingListHandler {
	return &ShoppingListHandler{repos: repos}
}

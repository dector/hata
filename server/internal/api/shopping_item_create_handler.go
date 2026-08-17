package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

// CreateItem handles POST /api/latest/house/{houseId}/shopping-list/{listId}/item
func (h *ShoppingListHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	_, houseID, ok := h.authorizeHouseAccess(w, r)
	if !ok {
		return
	}

	listUID := strings.TrimSpace(chi.URLParam(r, "listId"))
	if listUID == "" {
		WriteError(w, http.StatusBadRequest, "List ID is required", "invalid-request")
		return
	}

	list, err := h.repos.ShoppingList().GetByHouseAndUID(r.Context(), houseID, listUID)
	if err != nil {
		fmt.Printf("Error loading shopping list %q for house %q: %v\n", listUID, houseID, err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return
	}
	if list == nil {
		WriteError(w, http.StatusNotFound, "Shopping list not found", "not-found")
		return
	}

	var req ShoppingItemCreateRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body", "invalid-request")
		return
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		WriteError(w, http.StatusBadRequest, "Name is required", "invalid-request")
		return
	}

	itemUID, err := generateUID()
	if err != nil {
		fmt.Printf("Error generating shopping item uid: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return
	}

	item, err := h.repos.ShoppingItem().Create(r.Context(), houseID, listUID, itemUID, name)
	if err != nil {
		fmt.Printf("Error creating shopping item in list %q (house %q): %v\n", listUID, houseID, err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return
	}

	WriteJSON(w, http.StatusCreated, ShoppingItemResponse{Item: shoppingItemInfoFromData(item)})
}

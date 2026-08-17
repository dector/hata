package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

// UpdateItem handles PATCH /api/latest/house/{houseId}/shopping-list/{listId}/item/{itemId}
func (h *ShoppingListHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	_, houseID, ok := h.authorizeHouseAccess(w, r)
	if !ok {
		return
	}

	listUID := strings.TrimSpace(chi.URLParam(r, "listId"))
	if listUID == "" {
		WriteError(w, http.StatusBadRequest, "List ID is required", "invalid-request")
		return
	}
	itemUID := strings.TrimSpace(chi.URLParam(r, "itemId"))
	if itemUID == "" {
		WriteError(w, http.StatusBadRequest, "Item ID is required", "invalid-request")
		return
	}

	var req ShoppingItemUpdateRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body", "invalid-request")
		return
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		WriteError(w, http.StatusBadRequest, "Name is required", "invalid-request")
		return
	}

	item, err := h.repos.ShoppingItem().GetByHouseAndListAndUID(r.Context(), houseID, listUID, itemUID)
	if err != nil {
		fmt.Printf("Error loading shopping item %q in list %q (house %q): %v\n", itemUID, listUID, houseID, err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return
	}
	if item == nil || item.DeletedAt != nil {
		WriteError(w, http.StatusNotFound, "Shopping item not found", "not-found")
		return
	}

	if err := h.repos.ShoppingItem().UpdateName(r.Context(), houseID, listUID, itemUID, name); err != nil {
		fmt.Printf("Error updating shopping item name %q in list %q (house %q): %v\n", itemUID, listUID, houseID, err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return
	}

	updated, err := h.repos.ShoppingItem().GetByHouseAndListAndUID(r.Context(), houseID, listUID, itemUID)
	if err != nil {
		fmt.Printf("Error reloading shopping item %q in list %q (house %q): %v\n", itemUID, listUID, houseID, err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return
	}
	if updated == nil {
		WriteError(w, http.StatusNotFound, "Shopping item not found", "not-found")
		return
	}

	WriteJSON(w, http.StatusOK, ShoppingItemResponse{Item: shoppingItemInfoFromData(updated)})
}

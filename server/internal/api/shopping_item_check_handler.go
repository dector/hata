package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

// SetItemChecked handles PATCH /api/latest/house/{houseId}/shopping-list/{listId}/item/{itemId}/check
func (h *ShoppingListHandler) SetItemChecked(w http.ResponseWriter, r *http.Request) {
	userID, houseID, ok := h.authorizeHouseMutation(w, r)
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

	var req ShoppingItemCheckRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body", "invalid-request")
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

	var checkedAt *time.Time
	var checkedByUserID *int
	if req.Checked {
		now := time.Now()
		checkedAt = &now
		checkedByUserID = &userID
	}

	if err := h.repos.ShoppingItem().SetChecked(r.Context(), houseID, listUID, itemUID, checkedAt, checkedByUserID); err != nil {
		fmt.Printf("Error setting shopping item checked state %q in list %q (house %q): %v\n", itemUID, listUID, houseID, err)
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

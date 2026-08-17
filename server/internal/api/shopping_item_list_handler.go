package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

// ListItems handles GET /api/latest/house/{houseId}/shopping-list/{listId}/item
func (h *ShoppingListHandler) ListItems(w http.ResponseWriter, r *http.Request) {
	_, houseID, ok := h.authorizeHouseAccess(w, r)
	if !ok {
		return
	}

	listUID := strings.TrimSpace(chi.URLParam(r, "listId"))
	if listUID == "" {
		WriteError(w, http.StatusBadRequest, "List ID is required", "invalid-request")
		return
	}

	includeDeleted := false
	if raw := strings.TrimSpace(r.URL.Query().Get("include_deleted")); raw != "" {
		parsed, err := strconv.ParseBool(raw)
		if err != nil {
			WriteError(w, http.StatusBadRequest, "include_deleted must be true or false", "invalid-request")
			return
		}
		includeDeleted = parsed
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

	items, err := h.repos.ShoppingItem().ListByHouseAndListUID(r.Context(), houseID, listUID, includeDeleted)
	if err != nil {
		fmt.Printf("Error listing shopping items for list %q in house %q: %v\n", listUID, houseID, err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return
	}

	resp := ShoppingItemListResponse{Items: make([]ShoppingItemInfo, 0, len(items))}
	for _, item := range items {
		resp.Items = append(resp.Items, shoppingItemInfoFromData(item))
	}

	WriteJSON(w, http.StatusOK, resp)
}

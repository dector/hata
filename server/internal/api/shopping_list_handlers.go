package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

// ListByHouse handles GET /api/latest/house/{houseId}/shopping-list
func (h *ShoppingListHandler) ListByHouse(w http.ResponseWriter, r *http.Request) {
	userID, houseID, ok := h.authorizeHouseAccess(w, r)
	if !ok {
		return
	}
	_ = userID

	lists, err := h.repos.ShoppingList().ListByHouse(r.Context(), houseID)
	if err != nil {
		fmt.Printf("Error listing shopping lists for house %q: %v\n", houseID, err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return
	}

	resp := ShoppingListListResponse{ShoppingLists: make([]ShoppingListInfo, 0, len(lists))}
	for _, list := range lists {
		resp.ShoppingLists = append(resp.ShoppingLists, shoppingListInfoFromData(list))
	}

	WriteJSON(w, http.StatusOK, resp)
}

// GetByHouseAndUID handles GET /api/latest/house/{houseId}/shopping-list/{listId}
func (h *ShoppingListHandler) GetByHouseAndUID(w http.ResponseWriter, r *http.Request) {
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

	WriteJSON(w, http.StatusOK, ShoppingListGetResponse{ShoppingList: shoppingListInfoFromData(list)})
}

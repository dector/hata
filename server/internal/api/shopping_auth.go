package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

func (h *ShoppingListHandler) authorizeHouseAccess(w http.ResponseWriter, r *http.Request) (int, string, bool) {
	userID, err := userIDFromBearerToken(r, h.repos)
	if err != nil {
		if errors.Is(err, errUnauthorized) {
			WriteError(w, http.StatusUnauthorized, "Unauthorized", "unauthorized")
			return 0, "", false
		}
		fmt.Printf("Error validating token: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return 0, "", false
	}

	houseID := strings.TrimSpace(chi.URLParam(r, "houseId"))
	if houseID == "" {
		WriteError(w, http.StatusBadRequest, "House ID is required", "invalid-request")
		return 0, "", false
	}

	role, allowed, err := h.userCanAccessShoppingListHouse(r, userID, houseID)
	if err != nil {
		fmt.Printf("Error listing house roles: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return 0, "", false
	}
	if !allowed || role == "guest" {
		WriteError(w, http.StatusForbidden, "Forbidden", "forbidden")
		return 0, "", false
	}

	return userID, houseID, true
}

func (h *ShoppingListHandler) userCanAccessShoppingListHouse(r *http.Request, userID int, houseID string) (string, bool, error) {
	memberships, err := h.repos.HouseRole().ListByUser(r.Context(), userID)
	if err != nil {
		return "", false, err
	}

	for _, membership := range memberships {
		if membership.HouseID != houseID {
			continue
		}
		return strings.ToLower(strings.TrimSpace(membership.Role)), true, nil
	}

	return "", false, nil
}

package api

import (
	"errors"
	"fmt"
	"net/http"

	"hata/internal/db"
)

// HouseHandler provides house endpoints.
type HouseHandler struct {
	repos db.Repositories
}

// NewHouseHandler creates a new HouseHandler.
func NewHouseHandler(repos db.Repositories) *HouseHandler {
	return &HouseHandler{repos: repos}
}

// List handles GET /api/latest/house
func (h *HouseHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, err := userIDFromBearerToken(r, h.repos)
	if err != nil {
		if errors.Is(err, errUnauthorized) {
			WriteError(w, http.StatusUnauthorized, "Unauthorized", "unauthorized")
			return
		}
		fmt.Printf("Error validating token: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return
	}

	ctx := r.Context()
	memberships, err := h.repos.HouseRole().ListByUser(ctx, userID)
	if err != nil {
		fmt.Printf("Error listing houses: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return
	}

	resp := HouseListResponse{
		Houses: make([]HouseInfo, 0, len(memberships)),
	}
	for _, membership := range memberships {
		resp.Houses = append(resp.Houses, HouseInfo{
			ID:          membership.HouseID,
			DisplayName: membership.DisplayName,
			Role:        membership.Role,
		})
	}

	WriteJSON(w, http.StatusOK, resp)
}

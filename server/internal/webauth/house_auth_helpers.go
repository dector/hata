package webauth

import (
	"github.com/go-chi/chi/v5"
	"hata/internal/db"
	"net/http"
	"strings"
)

func (h *Handler) authorizedHouseMembership(w http.ResponseWriter, r *http.Request) *db.HouseMembershipData {
	if h.repos == nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return nil
	}
	auth, ok := AuthFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return nil
	}
	houseID := strings.TrimSpace(chi.URLParam(r, "houseId"))
	if houseID == "" {
		http.Error(w, "house id is required", http.StatusBadRequest)
		return nil
	}
	memberships, err := h.repos.HouseRole().ListByUser(r.Context(), auth.UserID)
	if err != nil {
		http.Error(w, "failed to load houses", http.StatusInternalServerError)
		return nil
	}
	membership := findMembership(memberships, houseID)
	if membership == nil {
		http.NotFound(w, r)
		return nil
	}
	return membership
}

func findMembership(memberships []*db.HouseMembershipData, houseID string) *db.HouseMembershipData {
	for _, membership := range memberships {
		if membership.HouseID == houseID {
			return membership
		}
	}
	return nil
}

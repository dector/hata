package webauth

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) DeleteHouseDiscoveryNetwork(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	membership := h.authorizedHouseMembership(w, r)
	if membership == nil {
		return
	}
	if !canManageHouse(membership.Role) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	networkID, err := strconv.Atoi(strings.TrimSpace(chi.URLParam(r, "networkId")))
	if err != nil || networkID <= 0 {
		http.Error(w, "invalid discovery network id", http.StatusBadRequest)
		return
	}
	if err := h.repos.HouseDiscoveryNetwork().DeleteByID(r.Context(), membership.HouseID, networkID); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	http.Redirect(w, r, houseManageURL(membership.HouseID), http.StatusSeeOther)
}

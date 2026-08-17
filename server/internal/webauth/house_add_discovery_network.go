package webauth

import (
	"net/http"
)

func (h *Handler) AddHouseDiscoveryNetwork(w http.ResponseWriter, r *http.Request) {
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
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}
	if _, err := h.repos.HouseDiscoveryNetwork().Create(r.Context(), membership.HouseID, r.FormValue("cidr"), r.FormValue("label")); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, houseManageURL(membership.HouseID), http.StatusSeeOther)
}

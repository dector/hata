package webauth

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) DeleteHouseDevice(w http.ResponseWriter, r *http.Request) {
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

	deviceID := strings.TrimSpace(chi.URLParam(r, "deviceId"))
	if deviceID == "" {
		http.Error(w, "device id is required", http.StatusBadRequest)
		return
	}
	if err := h.repos.Device().DeleteByHouseAndID(r.Context(), membership.HouseID, deviceID); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	http.Redirect(w, r, houseManageURL(membership.HouseID), http.StatusSeeOther)
}

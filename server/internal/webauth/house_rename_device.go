package webauth

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"
)

func (h *Handler) RenameHouseDevice(w http.ResponseWriter, r *http.Request) {
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
	deviceID := strings.TrimSpace(chi.URLParam(r, "deviceId"))
	if deviceID == "" {
		http.Error(w, "device id is required", http.StatusBadRequest)
		return
	}
	name := strings.TrimSpace(r.FormValue("name"))
	if name == "" {
		http.Error(w, "device name is required", http.StatusBadRequest)
		return
	}
	if err := h.repos.Device().UpdateName(r.Context(), membership.HouseID, deviceID, name); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	http.Redirect(w, r, houseManageURL(membership.HouseID), http.StatusSeeOther)
}

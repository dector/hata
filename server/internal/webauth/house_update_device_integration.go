package webauth

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) UpdateHouseDeviceIntegration(w http.ResponseWriter, r *http.Request) {
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
	device, err := h.repos.Device().GetByHouseAndID(r.Context(), membership.HouseID, deviceID)
	if err != nil {
		http.Error(w, "failed to load device", http.StatusInternalServerError)
		return
	}
	if device == nil {
		http.NotFound(w, r)
		return
	}

	integration := strings.ToLower(strings.TrimSpace(r.FormValue("integration")))
	if integration == "" {
		integration = strings.ToLower(strings.TrimSpace(device.IntegrationID))
	}
	if integration != strings.ToLower(strings.TrimSpace(device.IntegrationID)) || integration != "wiz" {
		http.Error(w, "unsupported integration", http.StatusBadRequest)
		return
	}

	ip := strings.TrimSpace(r.FormValue("ip"))
	if net.ParseIP(ip) == nil {
		http.Error(w, "valid device IP is required", http.StatusBadRequest)
		return
	}
	data := map[string]string{"ip": ip}
	if mac := normalizeMAC(r.FormValue("mac")); mac != "" {
		data["mac"] = mac
	} else if existingMAC := normalizeMAC(deviceIntegrationField(device, "mac")); existingMAC != "" {
		data["mac"] = existingMAC
	}
	buf, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "invalid integration data", http.StatusBadRequest)
		return
	}
	integrationData := string(buf)
	if err := h.repos.Device().UpdateIntegrationData(r.Context(), membership.HouseID, deviceID, &integrationData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	state := strings.ToLower(strings.TrimSpace(r.FormValue("state")))
	if state == "on" || state == "off" {
		_ = h.repos.Device().UpdateStatus(r.Context(), membership.HouseID, deviceID, state, "online")
	} else {
		_ = h.repos.Device().UpdateAvailability(r.Context(), membership.HouseID, deviceID, "online")
	}
	http.Redirect(w, r, houseManageURL(membership.HouseID), http.StatusSeeOther)
}

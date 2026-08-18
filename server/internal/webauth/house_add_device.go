package webauth

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"
)

func (h *Handler) AddHouseDevice(w http.ResponseWriter, r *http.Request) {
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

	integration := strings.ToLower(strings.TrimSpace(r.FormValue("integration")))
	if integration != "wiz" {
		http.Error(w, "unsupported integration", http.StatusBadRequest)
		return
	}
	ip := strings.TrimSpace(r.FormValue("ip"))
	if net.ParseIP(ip) == nil {
		http.Error(w, "valid device IP is required", http.StatusBadRequest)
		return
	}
	name := strings.TrimSpace(r.FormValue("name"))
	if name == "" {
		name = "WiZ " + ip
	}
	state := strings.ToLower(strings.TrimSpace(r.FormValue("state")))
	if state != "on" && state != "off" {
		state = ""
	}

	data := map[string]string{"ip": ip}
	if mac := normalizeMAC(r.FormValue("mac")); mac != "" {
		data["mac"] = mac
	}
	buf, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "invalid integration data", http.StatusBadRequest)
		return
	}
	integrationData := string(buf)
	if _, err := h.repos.Device().Create(r.Context(), membership.HouseID, deviceIDForIntegrationIP(integration, ip), name, integration, &integrationData, state); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, houseManageURL(membership.HouseID), http.StatusSeeOther)
}

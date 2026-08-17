package webauth

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"
)

type discoveredDeviceView struct {
	Integration string `json:"integration"`
	Name        string `json:"name"`
	IP          string `json:"ip"`
	State       string `json:"state"`
	InHouse     bool   `json:"inHouse"`
	AddURL      string `json:"addUrl"`
}

func (h *Handler) HouseDeviceDiscovery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if h.repos == nil || h.discoverDevices == nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	auth, ok := AuthFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	houseID := strings.TrimSpace(chi.URLParam(r, "houseId"))
	if houseID == "" {
		http.Error(w, "house id is required", http.StatusBadRequest)
		return
	}

	memberships, err := h.repos.HouseRole().ListByUser(r.Context(), auth.UserID)
	if err != nil {
		http.Error(w, "failed to load houses", http.StatusInternalServerError)
		return
	}
	if findMembership(memberships, houseID) == nil {
		http.NotFound(w, r)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	existingDevices, err := h.repos.Device().ListByHouse(r.Context(), houseID)
	if err != nil {
		http.Error(w, "failed to load devices", http.StatusInternalServerError)
		return
	}
	existingIPs := deviceIPsByIntegration(existingDevices)

	discoveryNetworks, err := h.repos.HouseDiscoveryNetwork().ListByHouse(r.Context(), houseID)
	if err != nil {
		http.Error(w, "failed to load discovery networks", http.StatusInternalServerError)
		return
	}
	cidrs := make([]string, 0, len(discoveryNetworks))
	for _, network := range discoveryNetworks {
		cidrs = append(cidrs, network.CIDR)
	}

	writeSSE(w, "progress", discoveryProgressMessage(discoveryNetworks))
	flusher.Flush()

	found := 0
	for device := range h.discoverDevices(r.Context(), cidrs) {
		found++
		view := discoveredDeviceView{
			Integration: device.Integration,
			Name:        device.Name,
			IP:          device.IP,
			State:       device.State,
			InHouse:     existingIPs[strings.ToLower(device.Integration)+":"+device.IP],
			AddURL:      houseManageURL(houseID) + "/devices",
		}
		payload, err := json.Marshal(view)
		if err != nil {
			continue
		}
		writeSSE(w, "device", string(payload))
		flusher.Flush()
	}

	if found == 0 {
		writeSSE(w, "done", "Discovery complete. No devices found.")
	} else {
		writeSSE(w, "done", fmt.Sprintf("Discovery complete. Found %d device(s).", found))
	}
	flusher.Flush()
}

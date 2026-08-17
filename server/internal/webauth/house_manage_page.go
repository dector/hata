package webauth

import (
	"github.com/go-chi/chi/v5"
	"hata/internal/webui"
	"net/http"
	"strings"
)

func (h *Handler) HouseManagePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if h.repos == nil {
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

	membership := findMembership(memberships, houseID)
	if membership == nil {
		http.NotFound(w, r)
		return
	}

	devices, err := h.repos.Device().ListByHouse(r.Context(), houseID)
	if err != nil {
		http.Error(w, "failed to load devices", http.StatusInternalServerError)
		return
	}

	viewDevices := make([]webui.AppDeviceData, 0, len(devices))
	for _, d := range devices {
		viewDevices = append(viewDevices, webui.AppDeviceData{
			HouseID:          d.HouseID,
			ID:               d.ID,
			Name:             d.Name,
			IntegrationID:    d.IntegrationID,
			State:            d.State,
			Availability:     d.Availability,
			IsLight:          deviceSupportsLight(d),
			LightBrightness:  d.LightBrightness,
			LightColorPreset: d.LightColorPreset,
		})
	}

	discoveryNetworks, err := h.repos.HouseDiscoveryNetwork().ListByHouse(r.Context(), houseID)
	if err != nil {
		http.Error(w, "failed to load discovery networks", http.StatusInternalServerError)
		return
	}
	viewNetworks := make([]webui.HouseDiscoveryNetworkData, 0, len(discoveryNetworks))
	for _, n := range discoveryNetworks {
		viewNetworks = append(viewNetworks, webui.HouseDiscoveryNetworkData{ID: n.ID, CIDR: n.CIDR, Label: n.Label})
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := webui.HouseManagePage(webui.HouseManagePageData{
		DisplayName:       displayNameFromAuth(auth),
		DiscoveryNetworks: viewNetworks,
		CanManageHouse:    canManageHouse(membership.Role),
		House: webui.AppHouseData{
			ID:          membership.HouseID,
			DisplayName: membership.DisplayName,
			Role:        membership.Role,
			Devices:     viewDevices,
		},
	}).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render house manage page", http.StatusInternalServerError)
		return
	}
}

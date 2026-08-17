package webauth

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"hata/internal/webui"
	"net"
	"net/http"
	"sort"
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
	sortMemberships(memberships)
	headerHouses := make([]webui.AppHouseData, 0, len(memberships))
	for _, m := range memberships {
		headerHouses = append(headerHouses, webui.AppHouseData{
			ID:          m.HouseID,
			DisplayName: m.DisplayName,
			Role:        m.Role,
		})
	}
	setActiveHouseCookie(w, membership.HouseID)

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
			IntegrationIP:    deviceIntegrationIP(d.IntegrationData),
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
		DisplayName:              displayNameFromAuth(auth),
		HeaderHouses:             headerHouses,
		ActiveHouseID:            membership.HouseID,
		DefaultDiscoveryNetworks: localIPv4CIDRs(),
		DiscoveryNetworks:        viewNetworks,
		CanManageHouse:           canManageHouse(membership.Role),
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

func deviceIntegrationIP(raw *string) string {
	if raw == nil {
		return ""
	}
	var data map[string]any
	if err := json.Unmarshal([]byte(*raw), &data); err != nil {
		return ""
	}
	ip, _ := data["ip"].(string)
	return strings.TrimSpace(ip)
}

func localIPv4CIDRs() []string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil
	}
	cidrs := make([]string, 0)
	seen := map[string]struct{}{}
	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			ip, ipNet, err := net.ParseCIDR(addr.String())
			if err != nil || ip.To4() == nil || ipNet == nil {
				continue
			}
			ipNet.IP = ip.Mask(ipNet.Mask)
			cidr := ipNet.String()
			if _, ok := seen[cidr]; ok {
				continue
			}
			seen[cidr] = struct{}{}
			cidrs = append(cidrs, cidr)
		}
	}
	sort.Strings(cidrs)
	return cidrs
}

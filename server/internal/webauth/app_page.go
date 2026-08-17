package webauth

import (
	"hata/internal/webui"
	"net/http"
	"sort"
)

func (h *Handler) AppPage(w http.ResponseWriter, r *http.Request) {
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

	memberships, err := h.repos.HouseRole().ListByUser(r.Context(), auth.UserID)
	if err != nil {
		http.Error(w, "failed to load houses", http.StatusInternalServerError)
		return
	}

	devices, err := h.repos.Device().ListByUser(r.Context(), auth.UserID)
	if err != nil {
		http.Error(w, "failed to load devices", http.StatusInternalServerError)
		return
	}

	devicesByHouse := make(map[string][]webui.AppDeviceData, len(memberships))
	for _, d := range devices {
		devicesByHouse[d.HouseID] = append(devicesByHouse[d.HouseID], appDeviceDataFromDB(d))
	}

	houses := make([]webui.AppHouseData, 0, len(memberships))
	for _, m := range memberships {
		houses = append(houses, webui.AppHouseData{
			ID:          m.HouseID,
			DisplayName: m.DisplayName,
			Role:        m.Role,
			Devices:     devicesByHouse[m.HouseID],
		})
	}

	sort.Slice(houses, func(i, j int) bool {
		if houses[i].DisplayName == houses[j].DisplayName {
			return houses[i].ID < houses[j].ID
		}
		return houses[i].DisplayName < houses[j].DisplayName
	})

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := webui.AppPage(webui.AppPageData{DisplayName: displayNameFromAuth(auth), Houses: houses}).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render app page", http.StatusInternalServerError)
		return
	}
}

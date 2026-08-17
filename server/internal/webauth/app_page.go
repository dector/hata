package webauth

import (
	"hata/internal/webui"
	"net/http"
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

	sortMemberships(memberships)
	activeMembership := activeHouseFromRequest(r, memberships)

	headerHouses := make([]webui.AppHouseData, 0, len(memberships))
	for _, m := range memberships {
		headerHouses = append(headerHouses, webui.AppHouseData{
			ID:          m.HouseID,
			DisplayName: m.DisplayName,
			Role:        m.Role,
		})
	}

	houses := make([]webui.AppHouseData, 0, 1)
	activeHouseID := ""
	if activeMembership != nil {
		activeHouseID = activeMembership.HouseID
		setActiveHouseCookie(w, activeHouseID)

		devices, err := h.repos.Device().ListByHouse(r.Context(), activeHouseID)
		if err != nil {
			http.Error(w, "failed to load devices", http.StatusInternalServerError)
			return
		}

		activeHouse := webui.AppHouseData{
			ID:          activeMembership.HouseID,
			DisplayName: activeMembership.DisplayName,
			Role:        activeMembership.Role,
			Devices:     make([]webui.AppDeviceData, 0, len(devices)),
		}
		for _, d := range devices {
			activeHouse.Devices = append(activeHouse.Devices, appDeviceDataFromDB(d))
		}
		houses = append(houses, activeHouse)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := webui.AppPage(webui.AppPageData{DisplayName: displayNameFromAuth(auth), Houses: houses, HeaderHouses: headerHouses, ActiveHouseID: activeHouseID}).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render app page", http.StatusInternalServerError)
		return
	}
}

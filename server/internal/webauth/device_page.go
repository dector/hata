package webauth

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"hata/internal/api"
	"hata/internal/db"
	"hata/internal/webui"
)

func (h *Handler) DevicePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	data, ok := h.devicePageData(w, r)
	if !ok {
		return
	}

	switch r.URL.Query().Get("reload") {
	case "ok":
		data.ReloadMessage = "Device info reloaded."
	case "offline":
		data.ReloadError = "Device did not answer. Marked offline."
	case "error":
		data.ReloadError = "Failed to reload device info."
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := webui.DevicePage(data).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render device page", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ReloadDeviceInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	data, ok := h.devicePageData(w, r)
	if !ok {
		return
	}
	if h.deviceController == nil {
		redirectAfterPost(w, r, data.Device.DetailsURL+"?reload=error")
		return
	}

	device := &db.DeviceData{
		ID:               data.Device.ID,
		Name:             data.Device.Name,
		IntegrationID:    data.Device.IntegrationID,
		IntegrationData:  &data.Device.IntegrationData,
		State:            data.Device.State,
		Availability:     data.Device.Availability,
		HouseID:          data.Device.HouseID,
		LightBrightness:  data.Device.LightBrightness,
		LightColorPreset: data.Device.LightColorPreset,
	}
	if data.Device.IntegrationData == "" {
		device.IntegrationData = nil
	}

	status, err := h.deviceController.GetStatus(r.Context(), device)
	if err != nil {
		if api.IsDeviceNoAck(err) {
			_ = h.repos.Device().UpdateAvailability(r.Context(), data.Device.HouseID, data.Device.ID, "offline")
			redirectAfterPost(w, r, data.Device.DetailsURL+"?reload=offline")
			return
		}
		redirectAfterPost(w, r, data.Device.DetailsURL+"?reload=error")
		return
	}
	if err := h.repos.Device().UpdateStatus(r.Context(), data.Device.HouseID, data.Device.ID, status.State, status.Availability); err != nil {
		redirectAfterPost(w, r, data.Device.DetailsURL+"?reload=error")
		return
	}
	if status.LightBrightness != nil || status.LightColorPreset != nil {
		if err := h.repos.Device().UpdateLight(r.Context(), data.Device.HouseID, data.Device.ID, status.LightBrightness, status.LightColorPreset); err != nil {
			redirectAfterPost(w, r, data.Device.DetailsURL+"?reload=error")
			return
		}
	}
	redirectAfterPost(w, r, data.Device.DetailsURL+"?reload=ok")
}

func (h *Handler) devicePageData(w http.ResponseWriter, r *http.Request) (webui.DevicePageData, bool) {
	if h.repos == nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return webui.DevicePageData{}, false
	}
	auth, ok := AuthFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return webui.DevicePageData{}, false
	}
	deviceID := strings.TrimSpace(chi.URLParam(r, "deviceId"))
	if deviceID == "" {
		http.Error(w, "device id is required", http.StatusBadRequest)
		return webui.DevicePageData{}, false
	}

	memberships, err := h.repos.HouseRole().ListByUser(r.Context(), auth.UserID)
	if err != nil {
		http.Error(w, "failed to load houses", http.StatusInternalServerError)
		return webui.DevicePageData{}, false
	}
	sortMemberships(memberships)
	headerHouses := make([]webui.AppHouseData, 0, len(memberships))
	for _, m := range memberships {
		headerHouses = append(headerHouses, webui.AppHouseData{ID: m.HouseID, DisplayName: m.DisplayName, Role: m.Role})
	}

	activeMembership := activeHouseFromRequest(r, memberships)
	device, membership, err := h.findUserDevice(r, auth.UserID, deviceID, activeMembership)
	if err != nil {
		http.Error(w, "failed to load device", http.StatusInternalServerError)
		return webui.DevicePageData{}, false
	}
	if device == nil || membership == nil {
		http.NotFound(w, r)
		return webui.DevicePageData{}, false
	}
	setActiveHouseCookie(w, membership.HouseID)

	return webui.DevicePageData{
		DisplayName:   displayNameFromAuth(auth),
		HeaderHouses:  headerHouses,
		ActiveHouseID: membership.HouseID,
		House:         webui.AppHouseData{ID: membership.HouseID, DisplayName: membership.DisplayName, Role: membership.Role},
		Device:        appDeviceDataFromDB(device),
		ReloadURL:     webui.DeviceReloadPath(device.ID),
	}, true
}

func (h *Handler) findUserDevice(r *http.Request, userID int, deviceID string, activeMembership *db.HouseMembershipData) (*db.DeviceData, *db.HouseMembershipData, error) {
	memberships, err := h.repos.HouseRole().ListByUser(r.Context(), userID)
	if err != nil {
		return nil, nil, err
	}
	if activeMembership != nil {
		device, err := h.repos.Device().GetByHouseAndID(r.Context(), activeMembership.HouseID, deviceID)
		if err != nil {
			return nil, nil, err
		}
		if device != nil {
			return device, activeMembership, nil
		}
	}
	for _, membership := range memberships {
		if activeMembership != nil && membership.HouseID == activeMembership.HouseID {
			continue
		}
		device, err := h.repos.Device().GetByHouseAndID(r.Context(), membership.HouseID, deviceID)
		if err != nil {
			return nil, nil, err
		}
		if device != nil {
			return device, membership, nil
		}
	}
	return nil, nil, nil
}

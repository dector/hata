package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"hata/internal/db"
)

func (h *DeviceHandler) userCanAccessHouse(r *http.Request, userID int, houseID string) (bool, error) {
	memberships, err := h.repos.HouseRole().ListByUser(r.Context(), userID)
	if err != nil {
		return false, err
	}
	for _, membership := range memberships {
		if membership.HouseID == houseID {
			return true, nil
		}
	}
	return false, nil
}

func normalizeDeviceState(raw string) (string, bool) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "on":
		return "on", true
	case "off":
		return "off", true
	default:
		return "", false
	}
}

func deviceInfoFromData(device *db.DeviceData) DeviceInfo {
	info := DeviceInfo{
		ID:   device.ID,
		Name: device.Name,
		Integration: DeviceIntegrationInfo{
			ID:   device.IntegrationID,
			Data: decodeIntegrationData(device.IntegrationData),
		},
		State:        device.State,
		Availability: device.Availability,
		Capabilities: capabilitiesFromData(device),
	}
	if info.Capabilities.Light {
		info.Light = lightInfoFromData(device)
	}
	return info
}

func deviceSupportsLight(device *db.DeviceData) bool {
	integrationID := strings.ToLower(strings.TrimSpace(device.IntegrationID))
	return strings.HasPrefix(integrationID, "wiz")
}

func capabilitiesFromData(device *db.DeviceData) DeviceCapabilitiesInfo {
	if deviceSupportsLight(device) {
		return DeviceCapabilitiesInfo{Light: true, Brightness: true, ColorPresets: true}
	}
	return DeviceCapabilitiesInfo{Light: false}
}

func lightInfoFromData(device *db.DeviceData) *DeviceLightInfo {
	return &DeviceLightInfo{Brightness: device.LightBrightness, ColorPreset: device.LightColorPreset}
}

func decodeIntegrationData(raw *string) map[string]any {
	if raw == nil {
		return nil
	}
	if strings.TrimSpace(*raw) == "" {
		return nil
	}

	var data map[string]any
	if err := json.Unmarshal([]byte(*raw), &data); err != nil {
		return nil
	}

	return data
}

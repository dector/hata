package webauth

import (
	"hata/internal/db"
	"hata/internal/light"
	"hata/internal/webui"
	"net/http"
	"strings"
)

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

func appDeviceDataFromDB(device *db.DeviceData) webui.AppDeviceData {
	return webui.AppDeviceData{
		HouseID:          device.HouseID,
		ID:               device.ID,
		Name:             device.Name,
		IntegrationID:    device.IntegrationID,
		State:            device.State,
		Availability:     device.Availability,
		ToggleURL:        webui.HouseDeviceTogglePath(device.HouseID, device.ID),
		StateURL:         webui.HouseDeviceStatePath(device.HouseID, device.ID),
		LightURL:         webui.HouseDeviceLightPath(device.HouseID, device.ID),
		IsLight:          deviceSupportsLight(device),
		LightBrightness:  device.LightBrightness,
		LightColorPreset: device.LightColorPreset,
		LightPresets:     webLightPresets(),
	}
}

func renderDeviceCard(w http.ResponseWriter, r *http.Request, device *db.DeviceData) {
	renderDeviceCardData(w, r, appDeviceDataFromDB(device))
}

func renderDeviceCardData(w http.ResponseWriter, r *http.Request, data webui.AppDeviceData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := webui.AppDeviceCard(data).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render device card", http.StatusInternalServerError)
	}
}

func renderDeviceCardWithTrigger(w http.ResponseWriter, r *http.Request, device *db.DeviceData, message string) {
	data := appDeviceDataFromDB(device)
	data.ToggleError = message
	renderDeviceCardData(w, r, data)
}

func deviceSupportsLight(device *db.DeviceData) bool {
	integrationID := strings.ToLower(strings.TrimSpace(device.IntegrationID))
	return strings.HasPrefix(integrationID, "wiz")
}

func webLightPresets() []webui.LightPresetData {
	presets := make([]webui.LightPresetData, 0, len(light.Presets))
	for _, preset := range light.Presets {
		presets = append(presets, webui.LightPresetData{ID: preset.ID, Label: preset.Label, Hex: preset.Hex})
	}
	return presets
}

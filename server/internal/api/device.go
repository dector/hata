package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"hata/internal/db"
	"hata/internal/light"

	"github.com/go-chi/chi/v5"
)

// DeviceHandler provides device endpoints.
type DeviceHandler struct {
	repos      db.Repositories
	controller DeviceController
}

// NewDeviceHandler creates a new DeviceHandler.
func NewDeviceHandler(repos db.Repositories) *DeviceHandler {
	return &DeviceHandler{
		repos:      repos,
		controller: NewRealDeviceController(),
	}
}

// NewDeviceHandlerWithController creates a new DeviceHandler with explicit device controller.
func NewDeviceHandlerWithController(repos db.Repositories, controller DeviceController) *DeviceHandler {
	if controller == nil {
		controller = NewRealDeviceController()
	}
	return &DeviceHandler{repos: repos, controller: controller}
}

// ListByHouse handles GET /api/latest/house/{houseId}/device
func (h *DeviceHandler) ListByHouse(w http.ResponseWriter, r *http.Request) {
	userID, err := userIDFromBearerToken(r, h.repos)
	if err != nil {
		if errors.Is(err, errUnauthorized) {
			WriteError(w, http.StatusUnauthorized, "Unauthorized", "unauthorized")
			return
		}
		fmt.Printf("Error validating token: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return
	}

	houseID := chi.URLParam(r, "houseId")
	if strings.TrimSpace(houseID) == "" {
		WriteError(w, http.StatusBadRequest, "House ID is required", "invalid-request")
		return
	}

	allowed, err := h.userCanAccessHouse(r, userID, houseID)
	if err != nil {
		fmt.Printf("Error listing house roles: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return
	}
	if !allowed {
		WriteError(w, http.StatusForbidden, "Forbidden", "forbidden")
		return
	}

	devices, err := h.repos.Device().ListByHouse(r.Context(), houseID)
	if err != nil {
		fmt.Printf("Error listing devices for house %q: %v\n", houseID, err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return
	}

	resp := DeviceListByHouseResponse{
		Devices: make([]DeviceInfo, 0, len(devices)),
	}
	for _, device := range devices {
		resp.Devices = append(resp.Devices, deviceInfoFromData(device))
	}

	WriteJSON(w, http.StatusOK, resp)
}

// ListByUser handles GET /api/latest/device
func (h *DeviceHandler) ListByUser(w http.ResponseWriter, r *http.Request) {
	userID, err := userIDFromBearerToken(r, h.repos)
	if err != nil {
		if errors.Is(err, errUnauthorized) {
			WriteError(w, http.StatusUnauthorized, "Unauthorized", "unauthorized")
			return
		}
		fmt.Printf("Error validating token: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return
	}

	devices, err := h.repos.Device().ListByUser(r.Context(), userID)
	if err != nil {
		fmt.Printf("Error listing devices for user %d: %v\n", userID, err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return
	}

	resp := DeviceListResponse{
		Devices: make([]DeviceInfoWithHouse, 0, len(devices)),
	}
	for _, device := range devices {
		resp.Devices = append(resp.Devices, DeviceInfoWithHouse{
			DeviceInfo: deviceInfoFromData(device),
			House:      HouseRef{ID: device.HouseID},
		})
	}

	WriteJSON(w, http.StatusOK, resp)
}

// SetState handles PATCH /api/latest/house/{houseId}/device/{deviceId}/state
func (h *DeviceHandler) SetState(w http.ResponseWriter, r *http.Request) {
	userID, err := userIDFromBearerToken(r, h.repos)
	if err != nil {
		if errors.Is(err, errUnauthorized) {
			WriteError(w, http.StatusUnauthorized, "Unauthorized", "unauthorized")
			return
		}
		fmt.Printf("Error validating token: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return
	}

	houseID := strings.TrimSpace(chi.URLParam(r, "houseId"))
	if houseID == "" {
		WriteError(w, http.StatusBadRequest, "House ID is required", "invalid-request")
		return
	}

	deviceID := strings.TrimSpace(chi.URLParam(r, "deviceId"))
	if deviceID == "" {
		WriteError(w, http.StatusBadRequest, "Device ID is required", "invalid-request")
		return
	}

	allowed, err := h.userCanAccessHouse(r, userID, houseID)
	if err != nil {
		fmt.Printf("Error listing house roles: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return
	}
	if !allowed {
		WriteError(w, http.StatusForbidden, "Forbidden", "forbidden")
		return
	}

	var req DeviceSetStateRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body", "invalid-request")
		return
	}

	newState, ok := normalizeDeviceState(req.State)
	if !ok {
		WriteError(w, http.StatusBadRequest, "State must be either 'on' or 'off'", "invalid-request")
		return
	}

	device, err := h.repos.Device().GetByHouseAndID(r.Context(), houseID, deviceID)
	if err != nil {
		fmt.Printf("Error loading device %q in house %q: %v\n", deviceID, houseID, err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return
	}
	if device == nil {
		WriteError(w, http.StatusNotFound, "Device not found", "not-found")
		return
	}

	if err := h.controller.SetState(r.Context(), device, newState); err != nil {
		if errors.Is(err, errDeviceNoAck) {
			_ = h.repos.Device().UpdateAvailability(r.Context(), houseID, deviceID, "offline")
		}
		if errors.Is(err, errUnsupportedIntegration) {
			WriteError(w, http.StatusUnprocessableEntity, "Device integration is not supported", "unsupported-integration")
			return
		}
		if errors.Is(err, errInvalidIntegrationData) {
			WriteError(w, http.StatusInternalServerError, "Device integration data is invalid", "invalid-integration-data")
			return
		}
		if errors.Is(err, errDeviceNoAck) {
			WriteError(w, http.StatusGatewayTimeout, "Device did not acknowledge command", "device-no-ack")
			return
		}
		fmt.Printf("Error controlling device %q in house %q: %v\n", deviceID, houseID, err)
		WriteError(w, http.StatusBadGateway, "Failed to control physical device", "device-control-failed")
		return
	}

	if err := h.repos.Device().UpdateStatus(r.Context(), houseID, deviceID, newState, "online"); err != nil {
		fmt.Printf("Error updating state for device %q in house %q: %v\n", deviceID, houseID, err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return
	}

	WriteJSON(w, http.StatusOK, DeviceSetStateResponse{
		DeviceID: deviceID,
		House:    HouseRef{ID: houseID},
		State:    newState,
	})
}

// SetLight handles PATCH /api/latest/house/{houseId}/device/{deviceId}/light
func (h *DeviceHandler) SetLight(w http.ResponseWriter, r *http.Request) {
	userID, err := userIDFromBearerToken(r, h.repos)
	if err != nil {
		if errors.Is(err, errUnauthorized) {
			WriteError(w, http.StatusUnauthorized, "Unauthorized", "unauthorized")
			return
		}
		fmt.Printf("Error validating token: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return
	}

	houseID := strings.TrimSpace(chi.URLParam(r, "houseId"))
	deviceID := strings.TrimSpace(chi.URLParam(r, "deviceId"))
	if houseID == "" || deviceID == "" {
		WriteError(w, http.StatusBadRequest, "House ID and device ID are required", "invalid-request")
		return
	}

	allowed, err := h.userCanAccessHouse(r, userID, houseID)
	if err != nil {
		fmt.Printf("Error listing house roles: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return
	}
	if !allowed {
		WriteError(w, http.StatusForbidden, "Forbidden", "forbidden")
		return
	}

	var req DeviceSetLightRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body", "invalid-request")
		return
	}
	if req.Brightness == nil && req.ColorPreset == nil {
		WriteError(w, http.StatusBadRequest, "At least one light field is required", "invalid-request")
		return
	}
	if req.Brightness != nil && !light.IsValidBrightness(*req.Brightness) {
		WriteError(w, http.StatusBadRequest, "Brightness must be 0..100 in 5% steps", "invalid-brightness")
		return
	}
	if req.ColorPreset != nil {
		preset := strings.TrimSpace(*req.ColorPreset)
		req.ColorPreset = &preset
		if !light.IsValidPreset(preset) {
			WriteError(w, http.StatusBadRequest, "Unsupported color preset", "invalid-color-preset")
			return
		}
	}

	device, err := h.repos.Device().GetByHouseAndID(r.Context(), houseID, deviceID)
	if err != nil {
		fmt.Printf("Error loading device %q in house %q: %v\n", deviceID, houseID, err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return
	}
	if device == nil {
		WriteError(w, http.StatusNotFound, "Device not found", "not-found")
		return
	}
	if !deviceSupportsLight(device) {
		WriteError(w, http.StatusUnprocessableEntity, "Device does not support light controls", "device-not-light")
		return
	}
	if device.Availability != "online" {
		WriteError(w, http.StatusConflict, "Device is offline", "device-offline")
		return
	}

	if err := h.controller.SetLight(r.Context(), device, req.Brightness, req.ColorPreset); err != nil {
		if errors.Is(err, errDeviceNoAck) {
			_ = h.repos.Device().UpdateAvailability(r.Context(), houseID, deviceID, "offline")
			WriteError(w, http.StatusGatewayTimeout, "Device did not acknowledge command", "device-no-ack")
			return
		}
		if errors.Is(err, errUnsupportedIntegration) {
			WriteError(w, http.StatusUnprocessableEntity, "Light control is not supported", "unsupported-light-control")
			return
		}
		if errors.Is(err, errInvalidIntegrationData) {
			WriteError(w, http.StatusInternalServerError, "Device integration data is invalid", "invalid-integration-data")
			return
		}
		fmt.Printf("Error setting light for device %q in house %q: %v\n", deviceID, houseID, err)
		WriteError(w, http.StatusBadGateway, "Failed to control physical device", "device-control-failed")
		return
	}

	if err := h.repos.Device().UpdateLight(r.Context(), houseID, deviceID, req.Brightness, req.ColorPreset); err != nil {
		fmt.Printf("Error updating light for device %q in house %q: %v\n", deviceID, houseID, err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return
	}
	if req.Brightness != nil {
		device.LightBrightness = req.Brightness
	}
	if req.ColorPreset != nil {
		device.LightColorPreset = req.ColorPreset
	}

	WriteJSON(w, http.StatusOK, DeviceSetLightResponse{
		DeviceID:     deviceID,
		House:        HouseRef{ID: houseID},
		State:        device.State,
		Availability: device.Availability,
		Light:        lightInfoFromData(device),
	})
}

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

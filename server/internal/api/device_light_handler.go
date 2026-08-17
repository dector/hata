package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"hata/internal/light"

	"github.com/go-chi/chi/v5"
)

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

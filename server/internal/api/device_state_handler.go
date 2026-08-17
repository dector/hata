package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

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

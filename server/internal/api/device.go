package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"hata/internal/db"
	dbrepo "hata/internal/db/repo"

	"github.com/go-chi/chi/v5"
)

// DeviceHandler provides device endpoints.
type DeviceHandler struct {
	repos db.Repositories
}

// NewDeviceHandler creates a new DeviceHandler.
func NewDeviceHandler(repos db.Repositories) *DeviceHandler {
	return &DeviceHandler{repos: repos}
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

	if err := h.repos.Device().UpdateState(r.Context(), houseID, deviceID, newState); err != nil {
		if errors.Is(err, dbrepo.ErrDeviceNotFound) {
			WriteError(w, http.StatusNotFound, "Device not found", "not-found")
			return
		}
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
	return DeviceInfo{
		ID:   device.ID,
		Name: device.Name,
		Integration: DeviceIntegrationInfo{
			ID:   device.IntegrationID,
			Data: decodeIntegrationData(device.IntegrationData),
		},
		State: device.State,
	}
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

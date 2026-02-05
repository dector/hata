package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"hata/internal/db"

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

	ctx := r.Context()
	memberships, err := h.repos.HouseRole().ListByUser(ctx, userID)
	if err != nil {
		fmt.Printf("Error listing house roles: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return
	}

	allowed := false
	for _, membership := range memberships {
		if membership.HouseID == houseID {
			allowed = true
			break
		}
	}

	if !allowed {
		WriteError(w, http.StatusForbidden, "Forbidden", "forbidden")
		return
	}

	devices, err := h.repos.Device().ListByHouse(ctx, houseID)
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

	ctx := r.Context()
	devices, err := h.repos.Device().ListByUser(ctx, userID)
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

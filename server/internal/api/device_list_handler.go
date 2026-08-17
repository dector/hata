package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

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

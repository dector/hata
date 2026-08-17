package webauth

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"hata/internal/api"
	"net/http"
	"strings"
)

func (h *Handler) ToggleHouseDevice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if h.repos == nil || h.deviceController == nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	auth, ok := AuthFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	houseID := strings.TrimSpace(chi.URLParam(r, "houseId"))
	deviceID := strings.TrimSpace(chi.URLParam(r, "deviceId"))
	if houseID == "" || deviceID == "" {
		http.Error(w, "house id and device id are required", http.StatusBadRequest)
		return
	}

	memberships, err := h.repos.HouseRole().ListByUser(r.Context(), auth.UserID)
	if err != nil {
		http.Error(w, "failed to load houses", http.StatusInternalServerError)
		return
	}
	if findMembership(memberships, houseID) == nil {
		http.NotFound(w, r)
		return
	}

	device, err := h.repos.Device().GetByHouseAndID(r.Context(), houseID, deviceID)
	if err != nil {
		http.Error(w, "failed to load device", http.StatusInternalServerError)
		return
	}
	if device == nil {
		http.NotFound(w, r)
		return
	}

	newState := "on"
	if strings.EqualFold(strings.TrimSpace(device.State), "on") {
		newState = "off"
	}

	if err := h.deviceController.SetState(r.Context(), device, newState); err != nil {
		if api.IsDeviceNoAck(err) {
			_ = h.repos.Device().UpdateAvailability(r.Context(), houseID, deviceID, "offline")
			device.Availability = "offline"
		}
		fmt.Printf("Error toggling device %q in house %q: %v\n", deviceID, houseID, err)
		if isPartialRequest(r) {
			renderDeviceCardWithTrigger(w, r, device, "Failed to toggle device.")
			return
		}
		http.Error(w, "failed to toggle device", http.StatusBadGateway)
		return
	}
	if err := h.repos.Device().UpdateStatus(r.Context(), houseID, deviceID, newState, "online"); err != nil {
		http.Error(w, "failed to update device", http.StatusInternalServerError)
		return
	}

	if isPartialRequest(r) {
		device.State = newState
		device.Availability = "online"
		renderDeviceCard(w, r, device)
		return
	}
	redirectAfterPost(w, r, "/app")
}

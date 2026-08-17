package webauth

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"hata/internal/api"
	"hata/internal/light"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) SetHouseDeviceLight(w http.ResponseWriter, r *http.Request) {
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
	if !deviceSupportsLight(device) {
		http.Error(w, "device does not support light controls", http.StatusUnprocessableEntity)
		return
	}
	if device.Availability != "online" {
		http.Error(w, "device is offline", http.StatusConflict)
		return
	}
	if err := r.ParseMultipartForm(32 << 20); err != nil && !errors.Is(err, http.ErrNotMultipart) {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	var brightness *int
	if _, ok := r.Form["brightness"]; ok {
		brightnessValue, err := strconv.Atoi(strings.TrimSpace(r.FormValue("brightness")))
		if err != nil || !light.IsValidBrightness(brightnessValue) {
			http.Error(w, "invalid brightness", http.StatusBadRequest)
			return
		}
		brightness = &brightnessValue
	}

	var colorPreset *string
	if _, ok := r.Form["colorPreset"]; ok {
		colorPresetValue := strings.TrimSpace(r.FormValue("colorPreset"))
		if !light.IsValidPreset(colorPresetValue) {
			http.Error(w, "invalid color preset", http.StatusBadRequest)
			return
		}
		colorPreset = &colorPresetValue
	}
	if brightness == nil && colorPreset == nil {
		http.Error(w, "brightness or color preset is required", http.StatusBadRequest)
		return
	}

	if err := h.deviceController.SetLight(r.Context(), device, brightness, colorPreset); err != nil {
		if api.IsDeviceNoAck(err) {
			_ = h.repos.Device().UpdateAvailability(r.Context(), houseID, deviceID, "offline")
		}
		fmt.Printf("Error setting light for device %q in house %q: %v\n", deviceID, houseID, err)
		http.Error(w, "failed to set light", http.StatusBadGateway)
		return
	}
	if err := h.repos.Device().UpdateLight(r.Context(), houseID, deviceID, brightness, colorPreset); err != nil {
		http.Error(w, "failed to update device", http.StatusInternalServerError)
		return
	}
	redirectAfterPost(w, r, "/app")
}

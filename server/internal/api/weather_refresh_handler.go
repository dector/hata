package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

// RefreshWeather handles POST /api/latest/house/{houseId}/weather/refresh.
func (h *HouseHandler) RefreshWeather(w http.ResponseWriter, r *http.Request) {
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

	if h.weather == nil {
		WriteError(w, http.StatusServiceUnavailable, "Weather is not configured", "weather-unavailable")
		return
	}
	if err := h.weather.RefreshHouse(r.Context(), houseID); err != nil {
		fmt.Printf("Error refreshing weather for house %q: %v\n", houseID, err)
		WriteError(w, http.StatusBadGateway, "Failed to refresh weather", "weather-refresh-failed")
		return
	}

	weatherInfo, err := h.weatherInfo(r.Context(), houseID)
	if err != nil {
		fmt.Printf("Error loading weather for house %q: %v\n", houseID, err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return
	}

	WriteJSON(w, http.StatusOK, WeatherRefreshResponse{Weather: weatherInfo})
}

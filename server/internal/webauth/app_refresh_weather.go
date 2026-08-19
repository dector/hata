package webauth

import (
	"net/http"

	"hata/internal/webui"
)

func (h *Handler) RefreshHouseWeather(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	membership := h.authorizedHouseMembership(w, r)
	if membership == nil {
		return
	}
	if h.weather == nil {
		http.Error(w, "weather is not configured", http.StatusServiceUnavailable)
		return
	}

	if err := h.weather.RefreshHouse(r.Context(), membership.HouseID); err != nil {
		http.Error(w, "failed to refresh weather", http.StatusBadGateway)
		return
	}

	weatherData := h.appWeatherData(r, membership.HouseID)
	if weatherData == nil {
		http.Error(w, "weather is not available", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := webui.AppWeatherCard(*weatherData).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render weather card", http.StatusInternalServerError)
		return
	}
}

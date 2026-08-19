package api

import (
	"context"
	"fmt"
	"net/http"

	"hata/internal/extension"
	"hata/internal/weather"
)

const weatherRefreshAction = "refresh"

type weatherHouseActionHandler struct {
	weather *weather.Plugin
}

// NewWeatherHouseActionHandler adapts the weather plugin to generic extension actions.
func NewWeatherHouseActionHandler(weatherPlugin *weather.Plugin) extension.HouseActionHandler {
	return &weatherHouseActionHandler{weather: weatherPlugin}
}

func (h *weatherHouseActionHandler) ExtensionID() string {
	return weather.ExtensionID
}

func (h *weatherHouseActionHandler) HandleHouseAction(ctx context.Context, houseID string, action string, r *http.Request) (extension.HouseActionResult, error) {
	if h.weather == nil {
		return extension.HouseActionResult{}, fmt.Errorf("weather is not configured")
	}
	if action != weatherRefreshAction {
		return extension.HouseActionResult{}, fmt.Errorf("unknown weather action %q", action)
	}
	if err := h.weather.RefreshHouse(ctx, houseID); err != nil {
		return extension.HouseActionResult{}, err
	}
	weatherInfo, err := weatherInfoForService(ctx, h.weather, houseID)
	if err != nil {
		return extension.HouseActionResult{}, err
	}
	return extension.HouseActionResult{
		Status: http.StatusOK,
		Data:   WeatherRefreshResponse{Weather: weatherInfo},
	}, nil
}

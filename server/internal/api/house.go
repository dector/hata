package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"hata/internal/db"
	"hata/internal/weather"
)

// HouseHandler provides house endpoints.
type HouseHandler struct {
	repos   db.Repositories
	weather *weather.Plugin
}

// NewHouseHandler creates a new HouseHandler.
func NewHouseHandler(repos db.Repositories) *HouseHandler {
	return &HouseHandler{repos: repos}
}

// NewHouseHandlerWithWeather creates a HouseHandler with weather plugin data.
func NewHouseHandlerWithWeather(repos db.Repositories, weatherPlugin *weather.Plugin) *HouseHandler {
	return &HouseHandler{repos: repos, weather: weatherPlugin}
}

// List handles GET /api/latest/house
func (h *HouseHandler) List(w http.ResponseWriter, r *http.Request) {
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
	memberships, err := h.repos.HouseRole().ListByUser(ctx, userID)
	if err != nil {
		fmt.Printf("Error listing houses: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return
	}

	resp := HouseListResponse{
		Houses: make([]HouseInfo, 0, len(memberships)),
	}
	for _, membership := range memberships {
		weatherInfo, err := h.weatherInfo(ctx, membership.HouseID)
		if err != nil {
			fmt.Printf("Error loading weather for house %q: %v\n", membership.HouseID, err)
			weatherInfo = &WeatherInfo{Status: string(weather.StatusError), Error: "failed to load weather"}
		}
		resp.Houses = append(resp.Houses, HouseInfo{
			ID:          membership.HouseID,
			DisplayName: membership.DisplayName,
			Location:    membership.Location,
			Role:        membership.Role,
			Weather:     weatherInfo,
		})
	}

	WriteJSON(w, http.StatusOK, resp)
}

func (h *HouseHandler) weatherInfo(ctx context.Context, houseID string) (*WeatherInfo, error) {
	if h.weather == nil {
		return &WeatherInfo{Status: string(weather.StatusDisabled)}, nil
	}
	status, err := h.weather.CurrentStatus(ctx, houseID)
	if err != nil {
		return nil, err
	}
	return weatherInfoFromStatus(status), nil
}

func weatherInfoFromStatus(status *weather.Status) *WeatherInfo {
	if status == nil {
		return &WeatherInfo{Status: string(weather.StatusDisabled)}
	}
	info := &WeatherInfo{
		Status:        string(status.Status),
		LocationLabel: status.LocationLabel,
		Error:         status.Error,
	}
	if status.Status == weather.StatusOK || status.Status == weather.StatusStale {
		temperature := status.Temperature
		conditionCode := status.ConditionCode
		info.Temperature = &temperature
		info.TemperatureUnit = status.TemperatureUnit
		info.ConditionCode = &conditionCode
		info.ConditionText = status.ConditionText
		info.ConditionIcon = status.ConditionIcon
		info.HumidityPercent = status.HumidityPercent
		info.WindSpeed = status.WindSpeed
		if status.WindSpeedUnit != nil {
			info.WindSpeedUnit = *status.WindSpeedUnit
		}
		info.ObservedAt = status.ObservedAt.Format(time.RFC3339)
		info.UpdatedAt = status.UpdatedAt.Format(time.RFC3339)
	}
	return info
}

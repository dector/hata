package weather

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"hata/internal/extension"
)

const weatherRefreshAction = "refresh"

// APIInfo contains current weather extension data for a house.
type APIInfo struct {
	Status          string   `json:"status"`
	LocationLabel   string   `json:"locationLabel,omitempty"`
	Temperature     *float64 `json:"temperature,omitempty"`
	TemperatureUnit string   `json:"temperatureUnit,omitempty"`
	ConditionCode   *int     `json:"conditionCode,omitempty"`
	ConditionText   string   `json:"conditionText,omitempty"`
	ConditionIcon   string   `json:"conditionIcon,omitempty"`
	HumidityPercent *float64 `json:"humidityPercent,omitempty"`
	WindSpeed       *float64 `json:"windSpeed,omitempty"`
	WindSpeedUnit   string   `json:"windSpeedUnit,omitempty"`
	ObservedAt      string   `json:"observedAt,omitempty"`
	UpdatedAt       string   `json:"updatedAt,omitempty"`
	Error           string   `json:"error,omitempty"`
}

// APIRefreshResponse represents the force weather refresh endpoint response.
type APIRefreshResponse struct {
	Weather *APIInfo `json:"weather"`
}

type apiHouseExtraProvider struct {
	weather *Plugin
}

// NewAPIHouseExtraProvider adapts weather status to generic API house extras.
func NewAPIHouseExtraProvider(weatherPlugin *Plugin) extension.HouseExtraProvider {
	return &apiHouseExtraProvider{weather: weatherPlugin}
}

func (p *apiHouseExtraProvider) ExtensionID() string {
	return ExtensionID
}

func (p *apiHouseExtraProvider) HouseExtra(ctx context.Context, houseID string) (any, error) {
	return apiInfoForService(ctx, p.weather, houseID)
}

type apiHouseActionHandler struct {
	weather *Plugin
}

// NewAPIHouseActionHandler adapts the weather plugin to generic API extension actions.
func NewAPIHouseActionHandler(weatherPlugin *Plugin) extension.HouseActionHandler {
	return &apiHouseActionHandler{weather: weatherPlugin}
}

func (h *apiHouseActionHandler) ExtensionID() string {
	return ExtensionID
}

func (h *apiHouseActionHandler) HandleHouseAction(ctx context.Context, houseID string, action string, r *http.Request) (extension.HouseActionResult, error) {
	if h.weather == nil {
		return extension.HouseActionResult{}, fmt.Errorf("weather is not configured")
	}
	if action != weatherRefreshAction {
		return extension.HouseActionResult{}, fmt.Errorf("unknown weather action %q", action)
	}
	if err := h.weather.RefreshHouse(ctx, houseID); err != nil {
		return extension.HouseActionResult{}, err
	}
	weatherInfo, err := apiInfoForService(ctx, h.weather, houseID)
	if err != nil {
		return extension.HouseActionResult{}, err
	}
	return extension.HouseActionResult{
		Status: http.StatusOK,
		Data:   APIRefreshResponse{Weather: weatherInfo},
	}, nil
}

func apiInfoForService(ctx context.Context, weatherPlugin *Plugin, houseID string) (*APIInfo, error) {
	if weatherPlugin == nil {
		return &APIInfo{Status: string(StatusDisabled)}, nil
	}
	status, err := weatherPlugin.CurrentStatus(ctx, houseID)
	if err != nil {
		return nil, err
	}
	return apiInfoFromStatus(status), nil
}

func apiInfoFromStatus(status *Status) *APIInfo {
	if status == nil {
		return &APIInfo{Status: string(StatusDisabled)}
	}
	info := &APIInfo{
		Status:        string(status.Status),
		LocationLabel: status.LocationLabel,
		Error:         status.Error,
	}
	if status.Status == StatusOK || status.Status == StatusStale {
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

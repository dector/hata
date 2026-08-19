package api

import (
	"context"
	"time"

	"hata/internal/extension"
	"hata/internal/weather"
)

type weatherHouseExtraProvider struct {
	weather *weather.Plugin
}

// NewWeatherHouseExtraProvider adapts weather status to generic house extras.
func NewWeatherHouseExtraProvider(weatherPlugin *weather.Plugin) extension.HouseExtraProvider {
	return &weatherHouseExtraProvider{weather: weatherPlugin}
}

func (p *weatherHouseExtraProvider) ExtensionID() string {
	return weather.ExtensionID
}

func (p *weatherHouseExtraProvider) HouseExtra(ctx context.Context, houseID string) (any, error) {
	return weatherInfoForService(ctx, p.weather, houseID)
}

func weatherInfoForService(ctx context.Context, weatherPlugin *weather.Plugin, houseID string) (*WeatherInfo, error) {
	if weatherPlugin == nil {
		return &WeatherInfo{Status: string(weather.StatusDisabled)}, nil
	}
	status, err := weatherPlugin.CurrentStatus(ctx, houseID)
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

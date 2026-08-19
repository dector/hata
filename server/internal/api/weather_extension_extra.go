package api

import (
	"context"

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

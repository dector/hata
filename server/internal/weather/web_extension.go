package weather

import (
	"context"
	"fmt"
	"net/http"

	"hata/internal/extension"
	"hata/internal/webui"

	"github.com/a-h/templ"
)

type webExtension struct {
	weather *Plugin
}

// NewWebHouseCardProvider adapts weather status to generic house cards.
func NewWebHouseCardProvider(weatherPlugin *Plugin) extension.HouseCardProvider {
	return &webExtension{weather: weatherPlugin}
}

// NewWebHouseActionHandler adapts weather browser actions to generic extension actions.
func NewWebHouseActionHandler(weatherPlugin *Plugin) extension.HouseWebActionHandler {
	return &webExtension{weather: weatherPlugin}
}

func (e *webExtension) ExtensionID() string {
	return ExtensionID
}

func (e *webExtension) HouseCard(ctx context.Context, houseID string, r *http.Request) (templ.Component, error) {
	weatherData := appWeatherData(ctx, e.weather, houseID)
	if weatherData == nil {
		return nil, nil
	}
	return webui.AppWeatherCard(*weatherData), nil
}

func (e *webExtension) HandleHouseWebAction(ctx context.Context, houseID string, action string, r *http.Request) (extension.HouseWebActionResult, error) {
	if e.weather == nil {
		return extension.HouseWebActionResult{}, fmt.Errorf("weather is not configured")
	}
	if action != weatherRefreshAction {
		return extension.HouseWebActionResult{}, fmt.Errorf("unknown weather action %q", action)
	}
	if err := e.weather.RefreshHouse(ctx, houseID); err != nil {
		return extension.HouseWebActionResult{}, err
	}
	weatherData := appWeatherData(ctx, e.weather, houseID)
	if weatherData == nil {
		return extension.HouseWebActionResult{}, fmt.Errorf("weather is not available")
	}
	return extension.HouseWebActionResult{
		Status:      http.StatusOK,
		ContentType: "text/html; charset=utf-8",
		Render:      webui.AppWeatherCard(*weatherData).Render,
	}, nil
}

func appWeatherData(ctx context.Context, weatherPlugin *Plugin, houseID string) *webui.AppWeatherData {
	if weatherPlugin == nil {
		return nil
	}
	status, err := weatherPlugin.CurrentStatus(ctx, houseID)
	if err != nil || status == nil || (status.Status != StatusOK && status.Status != StatusStale) {
		return nil
	}
	windUnit := ""
	if status.WindSpeedUnit != nil {
		windUnit = *status.WindSpeedUnit
	}
	return &webui.AppWeatherData{
		HouseID:         houseID,
		RefreshURL:      webui.HouseExtensionActionPath(houseID, ExtensionID, weatherRefreshAction),
		Status:          string(status.Status),
		LocationLabel:   status.LocationLabel,
		Temperature:     status.Temperature,
		TemperatureUnit: status.TemperatureUnit,
		ConditionText:   status.ConditionText,
		ConditionIcon:   status.ConditionIcon,
		HumidityPercent: status.HumidityPercent,
		WindSpeed:       status.WindSpeed,
		WindSpeedUnit:   windUnit,
		UpdatedAt:       status.UpdatedAt,
	}
}

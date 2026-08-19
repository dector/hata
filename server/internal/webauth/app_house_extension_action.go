package webauth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"hata/internal/extension"
	"hata/internal/weather"
	"hata/internal/webui"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
)

const weatherRefreshAction = "refresh"

// HandleHouseExtensionAction handles POST /h/{houseId}/extension/{extensionId}/actions/{action}.
func (h *Handler) HandleHouseExtensionAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	membership := h.authorizedHouseMembership(w, r)
	if membership == nil {
		return
	}
	extensionID := strings.TrimSpace(chi.URLParam(r, "extensionId"))
	action := strings.TrimSpace(chi.URLParam(r, "action"))
	if extensionID == "" || action == "" {
		http.Error(w, "extension ID and action are required", http.StatusBadRequest)
		return
	}

	handler, ok := h.extensions.HouseWebAction(extensionID)
	if !ok {
		http.Error(w, "extension action not found", http.StatusNotFound)
		return
	}
	result, err := handler.HandleHouseWebAction(r.Context(), membership.HouseID, action, r)
	if err != nil {
		http.Error(w, "failed to handle extension action", http.StatusBadGateway)
		return
	}
	status := result.Status
	if status == 0 {
		status = http.StatusOK
	}
	contentType := result.ContentType
	if contentType == "" {
		contentType = "text/html; charset=utf-8"
	}
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(status)
	if result.Render != nil {
		if err := result.Render(r.Context(), w); err != nil {
			fmt.Printf("Error rendering extension action %q/%q for house %q: %v\n", extensionID, action, membership.HouseID, err)
		}
	}
}

type weatherWebExtension struct {
	weather *weather.Plugin
}

// NewWeatherHouseCardProvider adapts weather status to generic house cards.
func NewWeatherHouseCardProvider(weatherPlugin *weather.Plugin) extension.HouseCardProvider {
	return &weatherWebExtension{weather: weatherPlugin}
}

// NewWeatherHouseWebActionHandler adapts weather browser actions to generic extension actions.
func NewWeatherHouseWebActionHandler(weatherPlugin *weather.Plugin) extension.HouseWebActionHandler {
	return &weatherWebExtension{weather: weatherPlugin}
}

func (e *weatherWebExtension) ExtensionID() string {
	return weather.ExtensionID
}

func (e *weatherWebExtension) HouseCard(ctx context.Context, houseID string, r *http.Request) (templ.Component, error) {
	weatherData := appWeatherData(ctx, e.weather, houseID)
	if weatherData == nil {
		return nil, nil
	}
	return webui.AppWeatherCard(*weatherData), nil
}

func (e *weatherWebExtension) HandleHouseWebAction(ctx context.Context, houseID string, action string, r *http.Request) (extension.HouseWebActionResult, error) {
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

func appWeatherData(ctx context.Context, weatherPlugin *weather.Plugin, houseID string) *webui.AppWeatherData {
	if weatherPlugin == nil {
		return nil
	}
	status, err := weatherPlugin.CurrentStatus(ctx, houseID)
	if err != nil || status == nil || (status.Status != weather.StatusOK && status.Status != weather.StatusStale) {
		return nil
	}
	windUnit := ""
	if status.WindSpeedUnit != nil {
		windUnit = *status.WindSpeedUnit
	}
	return &webui.AppWeatherData{
		HouseID:         houseID,
		RefreshURL:      webui.HouseExtensionActionPath(houseID, weather.ExtensionID, weatherRefreshAction),
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

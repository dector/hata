package weather

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"hata/internal/db"
)

const (
	PluginScope = "hata.weather.v1"
	EnabledKey  = "enabled"
)

// StatusName describes weather plugin status for a house.
type StatusName string

const (
	StatusDisabled      StatusName = "disabled"
	StatusNotConfigured StatusName = "not_configured"
	StatusPending       StatusName = "pending"
	StatusOK            StatusName = "ok"
	StatusStale         StatusName = "stale"
	StatusError         StatusName = "error"
)

// Status is the weather plugin view of a house.
type Status struct {
	Status          StatusName
	LocationLabel   string
	Temperature     float64
	TemperatureUnit string
	ConditionCode   int
	ConditionText   string
	ConditionIcon   string
	HumidityPercent *float64
	WindSpeed       *float64
	WindSpeedUnit   *string
	ObservedAt      time.Time
	UpdatedAt       time.Time
	Error           string
}

// Plugin contains simple weather plugin behavior.
type Plugin struct {
	Repos          db.Repositories
	Provider       Provider
	RefreshEvery   time.Duration
	StaleAfter     time.Duration
	RequestTimeout time.Duration
}

// NewPlugin creates a weather plugin.
func NewPlugin(repos db.Repositories, provider Provider) *Plugin {
	return &Plugin{
		Repos:          repos,
		Provider:       provider,
		RefreshEvery:   15 * time.Minute,
		StaleAfter:     30 * time.Minute,
		RequestTimeout: 10 * time.Second,
	}
}

// CurrentStatus returns current cached weather plugin status for a house.
func (p *Plugin) CurrentStatus(ctx context.Context, houseID string) (*Status, error) {
	enabled, err := p.IsEnabled(ctx, houseID)
	if err != nil {
		return nil, err
	}
	if !enabled {
		return &Status{Status: StatusDisabled}, nil
	}

	house, err := p.Repos.House().GetByID(ctx, houseID)
	if err != nil {
		return nil, err
	}
	if house == nil {
		return nil, fmt.Errorf("house %q not found", houseID)
	}
	location, err := ParseLocation(house.Location)
	if err != nil {
		return &Status{Status: StatusError, Error: err.Error()}, nil
	}
	if location == nil {
		return &Status{Status: StatusNotConfigured}, nil
	}

	cached, err := p.Repos.WeatherStatus().GetByHouse(ctx, houseID)
	if err != nil {
		return nil, err
	}
	if cached == nil {
		return &Status{Status: StatusPending, LocationLabel: location.Label}, nil
	}

	status := statusFromCache(cached)
	status.LocationLabel = location.Label
	if p.StaleAfter > 0 && time.Since(cached.UpdatedAt) > p.StaleAfter {
		status.Status = StatusStale
	}
	return status, nil
}

// RefreshAll refreshes all enabled and configured houses.
func (p *Plugin) RefreshAll(ctx context.Context) {
	houses, err := p.Repos.House().ListAll(ctx)
	if err != nil {
		log.Printf("weather: failed listing houses: %v", err)
		return
	}
	for _, house := range houses {
		if err := p.RefreshHouse(ctx, house.ID); err != nil {
			log.Printf("weather: failed refreshing house %q: %v", house.ID, err)
		}
	}
}

// RefreshHouse refreshes one house if weather is enabled and location is configured.
func (p *Plugin) RefreshHouse(ctx context.Context, houseID string) error {
	if p.Provider == nil {
		return fmt.Errorf("weather provider is not configured")
	}
	enabled, err := p.IsEnabled(ctx, houseID)
	if err != nil || !enabled {
		return err
	}
	house, err := p.Repos.House().GetByID(ctx, houseID)
	if err != nil || house == nil {
		return err
	}
	location, err := ParseLocation(house.Location)
	if err != nil || location == nil {
		return err
	}

	requestCtx := ctx
	cancel := func() {}
	if p.RequestTimeout > 0 {
		requestCtx, cancel = context.WithTimeout(ctx, p.RequestTimeout)
	}
	defer cancel()

	current, err := p.Provider.Current(requestCtx, CurrentRequest{Location: Location{Latitude: location.Latitude, Longitude: location.Longitude}})
	if err != nil {
		return err
	}
	windUnit := current.WindSpeedUnit
	_, err = p.Repos.WeatherStatus().UpsertByHouse(ctx, db.WeatherStatusData{
		HouseID:         houseID,
		Temperature:     current.Temperature,
		TemperatureUnit: current.TemperatureUnit,
		ConditionCode:   current.ConditionCode,
		ConditionText:   current.ConditionText,
		ConditionIcon:   current.ConditionIcon,
		HumidityPercent: current.HumidityPercent,
		WindSpeed:       current.WindSpeed,
		WindSpeedUnit:   &windUnit,
		ObservedAt:      current.ObservedAt,
		UpdatedAt:       time.Now().UTC(),
	})
	return err
}

// IsEnabled checks the weather enabled metasetting. Missing means false.
func (p *Plugin) IsEnabled(ctx context.Context, houseID string) (bool, error) {
	setting, err := p.Repos.HouseMetaSetting().Get(ctx, houseID, PluginScope, EnabledKey)
	if err != nil || setting == nil {
		return false, err
	}
	return strings.EqualFold(strings.TrimSpace(setting.Value), "true"), nil
}

// StartPoller starts lightweight background weather refresh.
func (p *Plugin) StartPoller(ctx context.Context) {
	go func() {
		p.RefreshAll(ctx)
		interval := p.RefreshEvery
		if interval <= 0 {
			interval = 15 * time.Minute
		}
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				p.RefreshAll(ctx)
			}
		}
	}()
}

// ParsedLocation is parsed house.location data.
type ParsedLocation struct {
	Latitude  float64
	Longitude float64
	Label     string
}

// ParseLocation parses latitude:longitude[:description].
func ParseLocation(value *string) (*ParsedLocation, error) {
	if value == nil || strings.TrimSpace(*value) == "" {
		return nil, nil
	}
	parts := strings.SplitN(strings.TrimSpace(*value), ":", 3)
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid location format")
	}
	lat, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return nil, fmt.Errorf("invalid latitude")
	}
	lon, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return nil, fmt.Errorf("invalid longitude")
	}
	req := CurrentRequest{Location: Location{Latitude: lat, Longitude: lon}}
	if err := req.Validate(); err != nil {
		return nil, err
	}
	label := ""
	if len(parts) == 3 {
		label = strings.TrimSpace(parts[2])
	}
	return &ParsedLocation{Latitude: lat, Longitude: lon, Label: label}, nil
}

func statusFromCache(cached *db.WeatherStatusData) *Status {
	return &Status{
		Status:          StatusOK,
		Temperature:     cached.Temperature,
		TemperatureUnit: cached.TemperatureUnit,
		ConditionCode:   cached.ConditionCode,
		ConditionText:   cached.ConditionText,
		ConditionIcon:   cached.ConditionIcon,
		HumidityPercent: cached.HumidityPercent,
		WindSpeed:       cached.WindSpeed,
		WindSpeedUnit:   cached.WindSpeedUnit,
		ObservedAt:      cached.ObservedAt,
		UpdatedAt:       cached.UpdatedAt,
	}
}

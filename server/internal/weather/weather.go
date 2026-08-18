package weather

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// UnitSystem describes units requested from a weather provider.
type UnitSystem string

const (
	UnitSystemMetric   UnitSystem = "metric"
	UnitSystemImperial UnitSystem = "imperial"
)

// Location identifies the place for a weather request.
type Location struct {
	Latitude  float64
	Longitude float64
}

// CurrentRequest contains options for fetching current weather.
type CurrentRequest struct {
	Location Location
	Units    UnitSystem
}

// Current describes provider-independent current weather.
type Current struct {
	Temperature     float64
	TemperatureUnit string
	ConditionCode   int
	ConditionText   string
	ConditionIcon   string
	HumidityPercent *float64
	WindSpeed       *float64
	WindSpeedUnit   string
	ObservedAt      time.Time
}

// Provider fetches weather from an external source.
type Provider interface {
	Current(ctx context.Context, req CurrentRequest) (*Current, error)
}

// Validate validates weather request values.
func (r CurrentRequest) Validate() error {
	if r.Location.Latitude < -90 || r.Location.Latitude > 90 {
		return fmt.Errorf("latitude must be between -90 and 90: %f", r.Location.Latitude)
	}
	if r.Location.Longitude < -180 || r.Location.Longitude > 180 {
		return fmt.Errorf("longitude must be between -180 and 180: %f", r.Location.Longitude)
	}
	if r.Units == "" {
		return nil
	}
	if r.Units != UnitSystemMetric && r.Units != UnitSystemImperial {
		return fmt.Errorf("unsupported unit system: %s", r.Units)
	}
	return nil
}

// Normalize fills defaults.
func (r CurrentRequest) Normalize() CurrentRequest {
	if r.Units == "" {
		r.Units = UnitSystemMetric
	}
	return r
}

var ErrMissingCurrentWeather = errors.New("missing current weather data")

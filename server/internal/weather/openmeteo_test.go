package weather

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestOpenMeteoProviderCurrentSuccessMetric(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("latitude"); got != "52.2297" {
			t.Fatalf("latitude query = %q", got)
		}
		if got := r.URL.Query().Get("longitude"); got != "21.0122" {
			t.Fatalf("longitude query = %q", got)
		}
		if got := r.URL.Query().Get("temperature_unit"); got != "celsius" {
			t.Fatalf("temperature_unit query = %q", got)
		}
		if got := r.URL.Query().Get("wind_speed_unit"); got != "kmh" {
			t.Fatalf("wind_speed_unit query = %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"current_units": {
				"temperature_2m": "°C",
				"wind_speed_10m": "km/h"
			},
			"current": {
				"time": "2026-08-18T12:00",
				"temperature_2m": 21.4,
				"relative_humidity_2m": 58,
				"weather_code": 2,
				"wind_speed_10m": 12.1
			}
		}`))
	}))
	defer server.Close()

	provider := NewOpenMeteoProvider(server.Client())
	provider.BaseURL = server.URL

	current, err := provider.Current(context.Background(), CurrentRequest{
		Location: Location{Latitude: 52.2297, Longitude: 21.0122},
	})
	if err != nil {
		t.Fatalf("Current returned error: %v", err)
	}

	if current.Temperature != 21.4 || current.TemperatureUnit != "°C" {
		t.Fatalf("temperature = %v %q", current.Temperature, current.TemperatureUnit)
	}
	if current.ConditionCode != 2 || current.ConditionText != "Partly cloudy" || current.ConditionIcon != "partly_cloudy" {
		t.Fatalf("condition = %d %q %q", current.ConditionCode, current.ConditionText, current.ConditionIcon)
	}
	if current.HumidityPercent == nil || *current.HumidityPercent != 58 {
		t.Fatalf("humidity = %v", current.HumidityPercent)
	}
	if current.WindSpeed == nil || *current.WindSpeed != 12.1 || current.WindSpeedUnit != "km/h" {
		t.Fatalf("wind = %v %q", current.WindSpeed, current.WindSpeedUnit)
	}
	wantTime := time.Date(2026, 8, 18, 12, 0, 0, 0, time.UTC)
	if !current.ObservedAt.Equal(wantTime) {
		t.Fatalf("observed at = %s, want %s", current.ObservedAt, wantTime)
	}
}

func TestOpenMeteoProviderCurrentSuccessImperial(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("temperature_unit"); got != "fahrenheit" {
			t.Fatalf("temperature_unit query = %q", got)
		}
		if got := r.URL.Query().Get("wind_speed_unit"); got != "mph" {
			t.Fatalf("wind_speed_unit query = %q", got)
		}
		_, _ = w.Write([]byte(`{
			"current_units": {"temperature_2m": "°F", "wind_speed_10m": "mp/h"},
			"current": {"time": "2026-08-18T12:00:00Z", "temperature_2m": 70.5, "weather_code": 0}
		}`))
	}))
	defer server.Close()

	provider := NewOpenMeteoProvider(server.Client())
	provider.BaseURL = server.URL

	current, err := provider.Current(context.Background(), CurrentRequest{
		Location: Location{Latitude: 52.2297, Longitude: 21.0122},
		Units:    UnitSystemImperial,
	})
	if err != nil {
		t.Fatalf("Current returned error: %v", err)
	}
	if current.TemperatureUnit != "°F" || current.WindSpeedUnit != "mp/h" {
		t.Fatalf("units = %q %q", current.TemperatureUnit, current.WindSpeedUnit)
	}
}

func TestOpenMeteoProviderCurrentNon2xx(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad", http.StatusBadGateway)
	}))
	defer server.Close()

	provider := NewOpenMeteoProvider(server.Client())
	provider.BaseURL = server.URL

	_, err := provider.Current(context.Background(), CurrentRequest{Location: Location{Latitude: 1, Longitude: 2}})
	if err == nil || !strings.Contains(err.Error(), "status 502") {
		t.Fatalf("error = %v", err)
	}
}

func TestOpenMeteoProviderCurrentInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{`))
	}))
	defer server.Close()

	provider := NewOpenMeteoProvider(server.Client())
	provider.BaseURL = server.URL

	_, err := provider.Current(context.Background(), CurrentRequest{Location: Location{Latitude: 1, Longitude: 2}})
	if err == nil || !strings.Contains(err.Error(), "decode open-meteo response") {
		t.Fatalf("error = %v", err)
	}
}

func TestOpenMeteoProviderCurrentMissingCurrent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"current": {"time": "2026-08-18T12:00"}}`))
	}))
	defer server.Close()

	provider := NewOpenMeteoProvider(server.Client())
	provider.BaseURL = server.URL

	_, err := provider.Current(context.Background(), CurrentRequest{Location: Location{Latitude: 1, Longitude: 2}})
	if !errors.Is(err, ErrMissingCurrentWeather) {
		t.Fatalf("error = %v", err)
	}
}

func TestCurrentRequestValidate(t *testing.T) {
	cases := []CurrentRequest{
		{Location: Location{Latitude: -91, Longitude: 0}},
		{Location: Location{Latitude: 0, Longitude: 181}},
		{Location: Location{Latitude: 0, Longitude: 0}, Units: UnitSystem("kelvin")},
	}
	for _, tc := range cases {
		if err := tc.Validate(); err == nil {
			t.Fatalf("Validate(%+v) returned nil", tc)
		}
	}
}

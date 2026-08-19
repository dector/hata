package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"hata/internal/weather"

	"github.com/go-chi/chi/v5"
)

type fakeWeatherProvider struct {
	calls int
}

func (f *fakeWeatherProvider) Current(ctx context.Context, req weather.CurrentRequest) (*weather.Current, error) {
	f.calls++
	humidity := 61.0
	wind := 4.2
	return &weather.Current{
		Temperature:     21.5,
		TemperatureUnit: "°C",
		ConditionCode:   1,
		ConditionText:   "Mostly clear",
		ConditionIcon:   "mostly_clear",
		HumidityPercent: &humidity,
		WindSpeed:       &wind,
		WindSpeedUnit:   "km/h",
		ObservedAt:      time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC),
	}, nil
}

func TestRefreshWeather_OK(t *testing.T) {
	repos, cleanup := setupAuthTest(t)
	defer cleanup()

	ctx := context.Background()
	user := createTestUser(t, repos, "user@example.com")
	token := createTestSession(t, repos, user.ID, strings.Repeat("w", 40))

	house, err := repos.House().Create(ctx, "house-1", "Main House")
	if err != nil {
		t.Fatalf("Failed to create house: %v", err)
	}
	if _, err := repos.HouseRole().Assign(ctx, house.ID, user.ID, "owner"); err != nil {
		t.Fatalf("Failed to assign house role: %v", err)
	}
	location := "52.52:13.405:Berlin"
	if err := repos.House().UpdateLocation(ctx, house.ID, &location); err != nil {
		t.Fatalf("Failed to set house location: %v", err)
	}
	if _, err := repos.HouseMetaSetting().Set(ctx, house.ID, weather.PluginScope, weather.EnabledKey, "true"); err != nil {
		t.Fatalf("Failed to enable weather: %v", err)
	}

	provider := &fakeWeatherProvider{}
	handler := NewHouseHandlerWithWeather(repos, weather.NewPlugin(repos, provider))

	router := chi.NewRouter()
	router.Post("/api/latest/house/{houseId}/weather/refresh", handler.RefreshWeather)

	req := httptest.NewRequest(http.MethodPost, "/api/latest/house/house-1/weather/refresh", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}
	if provider.calls != 1 {
		t.Fatalf("Expected provider to be called once, got %d", provider.calls)
	}

	var resp WeatherRefreshResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if resp.Weather == nil || resp.Weather.Status != string(weather.StatusOK) || resp.Weather.ConditionText != "Mostly clear" {
		t.Fatalf("Unexpected weather response: %+v", resp.Weather)
	}
}

func TestRefreshWeather_Forbidden(t *testing.T) {
	repos, cleanup := setupAuthTest(t)
	defer cleanup()

	ctx := context.Background()
	user := createTestUser(t, repos, "user@example.com")
	token := createTestSession(t, repos, user.ID, strings.Repeat("x", 40))
	if _, err := repos.House().Create(ctx, "house-1", "Main House"); err != nil {
		t.Fatalf("Failed to create house: %v", err)
	}

	handler := NewHouseHandlerWithWeather(repos, weather.NewPlugin(repos, &fakeWeatherProvider{}))
	router := chi.NewRouter()
	router.Post("/api/latest/house/{houseId}/weather/refresh", handler.RefreshWeather)

	req := httptest.NewRequest(http.MethodPost, "/api/latest/house/house-1/weather/refresh", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("Expected status 403, got %d", w.Code)
	}
}

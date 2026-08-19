package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"hata/internal/extension"
	"hata/internal/weather"
)

func TestHouseList_IncludesDefaultWeatherExtra(t *testing.T) {
	repos, cleanup := setupAuthTest(t)
	defer cleanup()
	restoreDefaultHouseExtras(t, []string{weather.ExtensionID})

	ctx := context.Background()
	user := createTestUser(t, repos, "user@example.com")
	token := createTestSession(t, repos, user.ID, strings.Repeat("h", 40))

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

	registry := extension.NewRegistry()
	if err := registry.RegisterHouseExtra(weather.NewAPIHouseExtraProvider(weather.NewPlugin(repos, nil))); err != nil {
		t.Fatalf("Failed to register weather extra: %v", err)
	}
	handler := NewHouseHandlerWithExtensions(repos, registry)
	req := httptest.NewRequest(http.MethodGet, "/api/latest/house", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler.List(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp HouseListResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if len(resp.Houses) != 1 {
		t.Fatalf("Expected 1 house, got %d", len(resp.Houses))
	}
	weatherExtra, ok := resp.Houses[0].Extras[weather.ExtensionID].(map[string]any)
	if !ok {
		t.Fatalf("Expected weather extra %q, got %#v", weather.ExtensionID, resp.Houses[0].Extras)
	}
	if weatherExtra["status"] != string(weather.StatusPending) {
		t.Fatalf("Expected pending weather status, got %#v", weatherExtra["status"])
	}
	if weatherExtra["locationLabel"] != "Berlin" {
		t.Fatalf("Expected Berlin location label, got %#v", weatherExtra["locationLabel"])
	}
}

func TestHouseList_OmitsWeatherExtraWhenNotDefaultEnriched(t *testing.T) {
	repos, cleanup := setupAuthTest(t)
	defer cleanup()
	restoreDefaultHouseExtras(t, nil)

	ctx := context.Background()
	user := createTestUser(t, repos, "user@example.com")
	token := createTestSession(t, repos, user.ID, strings.Repeat("i", 40))

	house, err := repos.House().Create(ctx, "house-1", "Main House")
	if err != nil {
		t.Fatalf("Failed to create house: %v", err)
	}
	if _, err := repos.HouseRole().Assign(ctx, house.ID, user.ID, "owner"); err != nil {
		t.Fatalf("Failed to assign house role: %v", err)
	}

	registry := extension.NewRegistry()
	if err := registry.RegisterHouseExtra(weather.NewAPIHouseExtraProvider(weather.NewPlugin(repos, nil))); err != nil {
		t.Fatalf("Failed to register weather extra: %v", err)
	}
	handler := NewHouseHandlerWithExtensions(repos, registry)
	req := httptest.NewRequest(http.MethodGet, "/api/latest/house", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler.List(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var raw struct {
		Houses []map[string]any `json:"houses"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &raw); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if len(raw.Houses) != 1 {
		t.Fatalf("Expected 1 house, got %d", len(raw.Houses))
	}
	if _, ok := raw.Houses[0]["extras"]; ok {
		t.Fatalf("Expected extras to be omitted, got %#v", raw.Houses[0]["extras"])
	}
	if _, ok := raw.Houses[0]["weather"]; ok {
		t.Fatalf("Expected legacy weather field to be omitted, got %#v", raw.Houses[0]["weather"])
	}
}

func restoreDefaultHouseExtras(t *testing.T, extras []string) {
	t.Helper()
	previous := DefaultHouseExtras
	DefaultHouseExtras = extras
	t.Cleanup(func() {
		DefaultHouseExtras = previous
	})
}

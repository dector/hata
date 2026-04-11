package apiui

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestPages_IncludeTokenWidgetAndCurlSnippet(t *testing.T) {
	tests := []struct {
		name     string
		handler  http.HandlerFunc
		path     string
		contains []string
	}{
		{
			name:    "index",
			handler: NewIndexHandler().Index,
			path:    "/apiui",
			contains: []string{
				"id=\"tokenStatus\"",
				"/apiui/latest/device",
				"/apiui/latest/ping",
			},
		},
		{
			name:    "login",
			handler: NewAuthUIHandler().LoginPage,
			path:    "/apiui/latest/auth/login",
			contains: []string{
				"/api/latest/auth/login",
				"curl -X POST",
				"id=\"tokenStatus\"",
			},
		},
		{
			name:    "house-list",
			handler: NewHouseUIHandler().ListPage,
			path:    "/apiui/latest/house",
			contains: []string{
				"/api/latest/house",
				"Authorization: Bearer $TOKEN",
				"id=\"tokenStatus\"",
			},
		},
		{
			name:    "device-list",
			handler: NewDeviceUIHandler().ListPage,
			path:    "/apiui/latest/device",
			contains: []string{
				"/api/latest/device",
				"Authorization: Bearer $TOKEN",
				"id=\"tokenStatus\"",
			},
		},
		{
			name:    "ping",
			handler: NewPingUIHandler().PingPage,
			path:    "/apiui/latest/ping",
			contains: []string{
				"/api/latest/ping",
				"curl http://localhost:4501/api/latest/ping",
				"id=\"tokenStatus\"",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			w := httptest.NewRecorder()
			tc.handler(w, req)

			if w.Code != http.StatusOK {
				t.Fatalf("expected status 200, got %d", w.Code)
			}

			body := w.Body.String()
			for _, part := range tc.contains {
				if !strings.Contains(body, part) {
					t.Fatalf("expected body to contain %q", part)
				}
			}
		})
	}
}

func TestHouseDevicePage_UsesHouseIDParam(t *testing.T) {
	handler := NewHouseUIHandler()
	router := chi.NewRouter()
	router.Get("/apiui/latest/house/{houseId}/device", handler.HouseDevicePage)

	req := httptest.NewRequest(http.MethodGet, "/apiui/latest/house/house-abc/device", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, `value="house-abc"`) {
		t.Fatalf("expected rendered houseId value in page")
	}
	if !strings.Contains(body, "/api/latest/house/{houseId}/device") {
		t.Fatalf("expected endpoint description in page")
	}
}

package webauth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"

	"hata/internal/api"
	"hata/internal/util"
)

func TestDevicePage_ShowsDeviceDetails(t *testing.T) {
	repos, cleanup := setupAuthWebTest(t)
	defer cleanup()

	passwordHash, err := util.HashPassword("secret")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	user, err := repos.User().Create(context.Background(), "user@example.com", passwordHash, "User")
	if err != nil {
		t.Fatalf("create user: %v", err)
	}
	if _, err := repos.House().Create(context.Background(), "H1", "Main Home"); err != nil {
		t.Fatalf("create house: %v", err)
	}
	if _, err := repos.HouseRole().Assign(context.Background(), "H1", user.ID, "owner"); err != nil {
		t.Fatalf("assign role: %v", err)
	}
	integrationData := `{"ip":"192.168.1.41","mac":"aa:bb"}`
	if _, err := repos.Device().Create(context.Background(), "H1", "wiz-192-168-1-41", "Kitchen", "wiz", &integrationData, "on"); err != nil {
		t.Fatalf("create device: %v", err)
	}

	h := NewHandler(api.NewAuthHandler(repos), repos)
	router := chi.NewRouter()
	router.Get("/d/{deviceId}", h.DevicePage)
	req := httptest.NewRequest(http.MethodGet, "/d/wiz-192-168-1-41", nil)
	req = req.WithContext(context.WithValue(req.Context(), authContextKey{}, AuthContext{UserID: user.ID, Username: user.Username, DisplayName: user.DisplayName}))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	body := w.Body.String()
	for _, expected := range []string{"Kitchen", "Device ID", "wiz-192-168-1-41", "Main Home", "wiz", "192.168.1.41", "aa:bb", "Reload info"} {
		if !strings.Contains(body, expected) {
			t.Fatalf("expected %q on page", expected)
		}
	}
}

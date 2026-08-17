package webauth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"hata/internal/api"
	"hata/internal/db"
)

func setupAuthWebTest(t *testing.T) (db.Repositories, func()) {
	t.Helper()

	ctx := context.Background()
	dbInst, err := db.OpenTestDB(ctx)
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}

	cleanup := func() {
		_ = dbInst.Close()
	}
	return dbInst.Repos(), cleanup
}

type fakeWebDeviceController struct {
	err error
}

func (f *fakeWebDeviceController) SetState(ctx context.Context, device *db.DeviceData, state string) error {
	return f.err
}

func (f *fakeWebDeviceController) SetLight(ctx context.Context, device *db.DeviceData, brightness *int, colorPreset *string) error {
	return nil
}

func (f *fakeWebDeviceController) GetStatus(ctx context.Context, device *db.DeviceData) (api.DeviceStatus, error) {
	return api.DeviceStatus{State: device.State, Availability: device.Availability}, nil
}

func assertLogoutRedirectAndCookieCleared(t *testing.T, w *httptest.ResponseRecorder) {
	t.Helper()

	if w.Code != http.StatusSeeOther {
		t.Fatalf("expected 303, got %d", w.Code)
	}
	if got := w.Header().Get("Location"); got != "/auth/login" {
		t.Fatalf("expected /auth/login redirect, got %q", got)
	}

	var authCookie *http.Cookie
	for _, c := range w.Result().Cookies() {
		if c.Name == api.AuthTokenCookieName {
			authCookie = c
			break
		}
	}
	if authCookie == nil {
		t.Fatalf("expected auth cookie to be cleared")
	}
	if authCookie.Value != "-" {
		t.Fatalf("expected cleared cookie value '-', got %q", authCookie.Value)
	}
	if !authCookie.Expires.Equal(time.Unix(0, 0)) {
		t.Fatalf("expected expires unix 0, got %v", authCookie.Expires)
	}
}

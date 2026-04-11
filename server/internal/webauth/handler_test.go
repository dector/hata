package webauth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"hata/internal/api"
	"hata/internal/db"
	"hata/internal/util"
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

func TestLoginPage_Renders(t *testing.T) {
	h := NewHandler(api.NewAuthHandler(nil), nil)

	req := httptest.NewRequest(http.MethodGet, "/auth/login?then=%2Fme", nil)
	w := httptest.NewRecorder()

	h.LoginPage(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "/js/htmx-2.0.4.min.js") {
		t.Fatalf("expected local htmx script in page")
	}
	if !strings.Contains(body, `name="then" value="/me"`) {
		t.Fatalf("expected hidden then field")
	}
}

func TestLogin_SetsCookieAndRedirects(t *testing.T) {
	repos, cleanup := setupAuthWebTest(t)
	defer cleanup()

	passwordHash, err := util.HashPassword("secret")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	if _, err := repos.User().Create(context.Background(), "user@example.com", passwordHash, "User"); err != nil {
		t.Fatalf("create user: %v", err)
	}

	h := NewHandler(api.NewAuthHandler(repos), repos)

	form := strings.NewReader("username=user@example.com&password=secret&then=%2Fme%3Fx%3D1")
	req := httptest.NewRequest(http.MethodPost, "/auth/login", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	h.Login(w, req)

	if w.Code != http.StatusSeeOther {
		t.Fatalf("expected 303, got %d", w.Code)
	}
	if got := w.Header().Get("Location"); got != "/me?x=1" {
		t.Fatalf("expected redirect to /me?x=1, got %q", got)
	}

	cookies := w.Result().Cookies()
	found := false
	for _, c := range cookies {
		if c.Name == api.AuthTokenCookieName {
			found = true
			if c.Value == "" || c.Value == "-" {
				t.Fatalf("expected real auth cookie value")
			}
		}
	}
	if !found {
		t.Fatalf("expected auth cookie to be set")
	}
}

func TestLogin_InvalidCredentials_NoCookie(t *testing.T) {
	repos, cleanup := setupAuthWebTest(t)
	defer cleanup()

	h := NewHandler(api.NewAuthHandler(repos), repos)

	form := strings.NewReader("username=nope&password=bad")
	req := httptest.NewRequest(http.MethodPost, "/auth/login", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	h.Login(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}

	for _, c := range w.Result().Cookies() {
		if c.Name == api.AuthTokenCookieName {
			t.Fatalf("did not expect auth cookie on failure")
		}
	}
}

func TestLogout_ClearsCookieAndRedirects(t *testing.T) {
	h := NewHandler(api.NewAuthHandler(nil), nil)

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	w := httptest.NewRecorder()

	h.Logout(w, req)

	assertLogoutRedirectAndCookieCleared(t, w)
}

func TestLogoutPage_AutoSubmitsOnGet(t *testing.T) {
	h := NewHandler(api.NewAuthHandler(nil), nil)

	req := httptest.NewRequest(http.MethodGet, "/auth/logout", nil)
	w := httptest.NewRecorder()

	h.LogoutPage(w, req)

	assertLogoutRedirectAndCookieCleared(t, w)
}

func TestAppPage_GroupsDevicesByHouse(t *testing.T) {
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
		t.Fatalf("create house H1: %v", err)
	}
	if _, err := repos.House().Create(context.Background(), "H2", "Garage"); err != nil {
		t.Fatalf("create house H2: %v", err)
	}

	if _, err := repos.HouseRole().Assign(context.Background(), "H1", user.ID, "owner"); err != nil {
		t.Fatalf("assign house role H1: %v", err)
	}
	if _, err := repos.HouseRole().Assign(context.Background(), "H2", user.ID, "guest"); err != nil {
		t.Fatalf("assign house role H2: %v", err)
	}

	if _, err := repos.Device().Create(context.Background(), "H1", "lamp-1", "Bedroom Lamp", "dummy", nil, "on"); err != nil {
		t.Fatalf("create device: %v", err)
	}

	h := NewHandler(api.NewAuthHandler(repos), repos)
	req := httptest.NewRequest(http.MethodGet, "/app", nil)
	req = req.WithContext(context.WithValue(req.Context(), authContextKey{}, AuthContext{UserID: user.ID, Username: user.Username}))
	w := httptest.NewRecorder()

	h.AppPage(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "Main Home") || !strings.Contains(body, "Garage") {
		t.Fatalf("expected both houses on page")
	}
	if !strings.Contains(body, "Bedroom Lamp") {
		t.Fatalf("expected house devices on page")
	}
	if !strings.Contains(body, "No devices in this house") {
		t.Fatalf("expected empty house message")
	}
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

func TestSanitizeThen(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "empty", in: "", want: "/me"},
		{name: "relative", in: "/me?a=1", want: "/me?a=1"},
		{name: "double-slash", in: "//evil.com", want: "/me"},
		{name: "absolute", in: "https://evil.com", want: "/me"},
		{name: "malformed-escape", in: "/me?x=%zz", want: "/me"},
		{name: "backslash", in: "/\\evil", want: "/me"},
		{name: "newline", in: "/me\nfoo", want: "/me"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sanitizeThen(tt.in); got != tt.want {
				t.Fatalf("sanitizeThen(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

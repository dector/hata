package webauth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"hata/internal/api"
	"hata/internal/util"
)

func TestLoginPage_Renders(t *testing.T) {
	h := NewHandler(api.NewAuthHandler(nil), nil)

	req := httptest.NewRequest(http.MethodGet, "/auth/login?then=%2Fme", nil)
	w := httptest.NewRecorder()

	h.LoginPage(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	body := w.Body.String()
	if strings.Contains(body, "htmx") || strings.Contains(body, "htmx-2.0.4.min.js") {
		t.Fatalf("did not expect htmx script in page")
	}
	if !strings.Contains(body, "/assets/js/datastar-1.0.2.js") {
		t.Fatalf("expected local datastar script in page")
	}
	if !strings.Contains(body, `data-on-submit__prevent="@post('/auth/login', {contentType: 'form', selector: '#login-shell'})"`) {
		t.Fatalf("expected datastar login submit handler")
	}
	if strings.Contains(body, "hx-post") {
		t.Fatalf("did not expect htmx form attributes")
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

func TestLogin_DatastarRedirectsWithJavaScript(t *testing.T) {
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
	req.Header.Set("Datastar-Request", "true")
	w := httptest.NewRecorder()

	h.Login(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if contentType := w.Header().Get("Content-Type"); !strings.Contains(contentType, "text/javascript") {
		t.Fatalf("expected javascript content type, got %q", contentType)
	}
	if got := w.Body.String(); got != `window.location.assign("/me?x=1");` {
		t.Fatalf("unexpected redirect script: %q", got)
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

func TestLogin_InvalidCredentials_DatastarReturnsPatchableHTML(t *testing.T) {
	repos, cleanup := setupAuthWebTest(t)
	defer cleanup()

	h := NewHandler(api.NewAuthHandler(repos), repos)

	form := strings.NewReader("username=nope&password=bad")
	req := httptest.NewRequest(http.MethodPost, "/auth/login", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Datastar-Request", "true")
	w := httptest.NewRecorder()

	h.Login(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, `id="login-shell"`) || !strings.Contains(body, "Invalid username or password") {
		t.Fatalf("expected patchable login shell with error, got %q", body)
	}
}

func TestLogout_ClearsCookieAndRedirects(t *testing.T) {
	h := NewHandler(api.NewAuthHandler(nil), nil)

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	w := httptest.NewRecorder()

	h.Logout(w, req)

	assertLogoutRedirectAndCookieCleared(t, w)
}

func TestLogout_DatastarRedirectsWithJavaScript(t *testing.T) {
	h := NewHandler(api.NewAuthHandler(nil), nil)

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	req.Header.Set("Datastar-Request", "true")
	w := httptest.NewRecorder()

	h.Logout(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if got := w.Body.String(); got != `window.location.assign("/auth/login");` {
		t.Fatalf("unexpected redirect script: %q", got)
	}
}

func TestLogoutPage_AutoSubmitsOnGet(t *testing.T) {
	h := NewHandler(api.NewAuthHandler(nil), nil)

	req := httptest.NewRequest(http.MethodGet, "/auth/logout", nil)
	w := httptest.NewRecorder()

	h.LogoutPage(w, req)

	assertLogoutRedirectAndCookieCleared(t, w)
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

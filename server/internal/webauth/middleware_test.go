package webauth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"hata/internal/api"
	"hata/internal/util"
)

func TestRequirePageAuth_RedirectsAnonymousUser(t *testing.T) {
	repos, cleanup := setupAuthWebTest(t)
	defer cleanup()

	mw := RequirePageAuth(repos)

	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusSeeOther {
		t.Fatalf("expected 303, got %d", w.Code)
	}
	if got := w.Header().Get("Location"); got != "/auth/login?then=%2Fme" {
		t.Fatalf("expected redirect to login with then=/me, got %q", got)
	}
}

func TestRequirePageAuth_RedirectsAnonymousUser_WithQuery(t *testing.T) {
	repos, cleanup := setupAuthWebTest(t)
	defer cleanup()

	mw := RequirePageAuth(repos)
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/me?tab=devices", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusSeeOther {
		t.Fatalf("expected 303, got %d", w.Code)
	}
	if got := w.Header().Get("Location"); got != "/auth/login?then=%2Fme%3Ftab%3Ddevices" {
		t.Fatalf("expected redirect with encoded query in then, got %q", got)
	}
}

func TestRequirePageAuth_AllowsValidCookie(t *testing.T) {
	repos, cleanup := setupAuthWebTest(t)
	defer cleanup()

	ctx := context.Background()
	passwordHash, err := util.HashPassword("secret")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	user, err := repos.User().Create(ctx, "u@example.com", passwordHash, "User")
	if err != nil {
		t.Fatalf("create user: %v", err)
	}
	session, err := repos.Session().Create(ctx, user.ID, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", time.Now().Add(24*time.Hour))
	if err != nil {
		t.Fatalf("create session: %v", err)
	}

	mw := RequirePageAuth(repos)

	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth, ok := AuthFromContext(r.Context())
		if !ok {
			t.Fatalf("expected auth context")
		}
		if auth.UserID != user.ID {
			t.Fatalf("expected user id %d, got %d", user.ID, auth.UserID)
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.AddCookie(&http.Cookie{Name: api.AuthTokenCookieName, Value: session.Token})
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

package webauth

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"hata/internal/api"
	"hata/internal/util"

	"github.com/go-chi/chi/v5"
)

func TestAppPage_ShowsOnlyActiveHouse(t *testing.T) {
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
	if _, err := repos.Device().Create(context.Background(), "H1", "lamp-2", "Desk Lamp", "dummy", nil, "off"); err != nil {
		t.Fatalf("create offline device: %v", err)
	}
	if err := repos.Device().UpdateAvailability(context.Background(), "H1", "lamp-2", "offline"); err != nil {
		t.Fatalf("mark device offline: %v", err)
	}

	h := NewHandler(api.NewAuthHandler(repos), repos)
	req := httptest.NewRequest(http.MethodGet, "/app", nil)
	req.AddCookie(&http.Cookie{Name: activeHouseCookieName, Value: "H1"})
	req = req.WithContext(context.WithValue(req.Context(), authContextKey{}, AuthContext{UserID: user.ID, Username: user.Username}))
	w := httptest.NewRecorder()

	h.AppPage(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "Main Home") || !strings.Contains(body, "Garage") {
		t.Fatalf("expected active house and switcher house on page")
	}
	if !strings.Contains(body, `href="/h/H1/manage"`) || !strings.Contains(body, "Manage house") {
		t.Fatalf("expected active house manage action in selector")
	}
	if strings.Contains(body, `href="/h/H2/manage"`) {
		t.Fatalf("did not expect inactive house manage link")
	}
	if !strings.Contains(body, "Bedroom Lamp") || !strings.Contains(body, "Desk Lamp") {
		t.Fatalf("expected active house devices on page")
	}
	if !strings.Contains(body, "offline") {
		t.Fatalf("expected offline device status on page")
	}
	if strings.Contains(body, "No devices in this house") {
		t.Fatalf("did not expect inactive house empty message")
	}
	if !strings.Contains(body, `id="device-card-H1-lamp-1"`) {
		t.Fatalf("expected stable device card id")
	}
	if !strings.Contains(body, `data-on-submit__prevent="@post(&#34;/h/H1/device/lamp-1/toggle&#34;, {contentType: &#39;form&#39;, selector: &#34;#device-card-H1-lamp-1&#34;})"`) {
		t.Fatalf("expected datastar device toggle handler")
	}
	if strings.Contains(body, "hx-post") || strings.Contains(body, "htmx:beforeRequest") || strings.Contains(body, "htmx-request") {
		t.Fatalf("did not expect htmx device toggle usage")
	}
}

func TestSetActiveHouse_SetsCookieAndRedirects(t *testing.T) {
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
	if _, err := repos.HouseRole().Assign(context.Background(), "H1", user.ID, "owner"); err != nil {
		t.Fatalf("assign house role H1: %v", err)
	}

	h := NewHandler(api.NewAuthHandler(repos), repos)
	req := httptest.NewRequest(http.MethodPost, "/app/active-house", strings.NewReader("house_id=H1&then=/app?x=1"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = req.WithContext(context.WithValue(req.Context(), authContextKey{}, AuthContext{UserID: user.ID, Username: user.Username}))
	w := httptest.NewRecorder()

	h.SetActiveHouse(w, req)

	if w.Code != http.StatusSeeOther {
		t.Fatalf("expected 303, got %d", w.Code)
	}
	if got := w.Header().Get("Location"); got != "/app?x=1" {
		t.Fatalf("expected redirect to then path, got %q", got)
	}
	var activeCookie *http.Cookie
	for _, c := range w.Result().Cookies() {
		if c.Name == activeHouseCookieName {
			activeCookie = c
			break
		}
	}
	if activeCookie == nil || activeCookie.Value != "H1" || !activeCookie.HttpOnly || activeCookie.SameSite != http.SameSiteLaxMode {
		t.Fatalf("expected secure active house cookie, got %#v", activeCookie)
	}
}

func TestSetActiveHouse_RejectsUnauthorizedHouse(t *testing.T) {
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

	h := NewHandler(api.NewAuthHandler(repos), repos)
	req := httptest.NewRequest(http.MethodPost, "/app/active-house", strings.NewReader("house_id=H1"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = req.WithContext(context.WithValue(req.Context(), authContextKey{}, AuthContext{UserID: user.ID, Username: user.Username}))
	w := httptest.NewRecorder()

	h.SetActiveHouse(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
	for _, c := range w.Result().Cookies() {
		if c.Name == activeHouseCookieName {
			t.Fatalf("did not expect active house cookie")
		}
	}
}

func TestPartialRequestDetection(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	if isPartialRequest(req) {
		t.Fatalf("plain request should not be partial")
	}

	req.Header.Set("Datastar-Request", "true")
	if !isPartialRequest(req) {
		t.Fatalf("datastar request should be partial")
	}

	req.Header.Del("Datastar-Request")
	req.Header.Set("HX-Request", "true")
	if isPartialRequest(req) {
		t.Fatalf("htmx request should not be partial")
	}
}

func TestToggleHouseDevice_DatastarUpdatesOnlyDeviceCard(t *testing.T) {
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
		t.Fatalf("assign house role: %v", err)
	}
	if _, err := repos.Device().Create(context.Background(), "H1", "lamp-1", "Bedroom Lamp", "dummy", nil, "off"); err != nil {
		t.Fatalf("create device: %v", err)
	}

	h := NewHandler(api.NewAuthHandler(repos), repos)
	h.deviceController = &fakeWebDeviceController{}
	router := chi.NewRouter()
	router.Post("/h/{houseId}/device/{deviceId}/toggle", h.ToggleHouseDevice)

	req := httptest.NewRequest(http.MethodPost, "/h/H1/device/lamp-1/toggle", nil)
	req.Header.Set("Datastar-Request", "true")
	req = req.WithContext(context.WithValue(req.Context(), authContextKey{}, AuthContext{UserID: user.ID, Username: user.Username}))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	body := w.Body.String()
	if !strings.Contains(body, `class="device-card`) || !strings.Contains(body, `data-state="on"`) {
		t.Fatalf("expected updated device card, got %s", body)
	}
	if !strings.Contains(body, `id="device-card-H1-lamp-1"`) {
		t.Fatalf("expected stable device card id, got %s", body)
	}
	if strings.Contains(body, "My devices") || strings.Contains(body, "Main Home") {
		t.Fatalf("expected only card partial, got full page: %s", body)
	}
	if strings.Contains(w.Header().Get("HX-Redirect"), "/app") {
		t.Fatalf("did not expect HX redirect")
	}
}

func TestToggleHouseDevice_DatastarFailureReturnsCardWithError(t *testing.T) {
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
		t.Fatalf("assign house role: %v", err)
	}
	if _, err := repos.Device().Create(context.Background(), "H1", "lamp-1", "Bedroom Lamp", "dummy", nil, "off"); err != nil {
		t.Fatalf("create device: %v", err)
	}

	h := NewHandler(api.NewAuthHandler(repos), repos)
	h.deviceController = &fakeWebDeviceController{err: errors.New("toggle failed")}
	router := chi.NewRouter()
	router.Post("/h/{houseId}/device/{deviceId}/toggle", h.ToggleHouseDevice)

	req := httptest.NewRequest(http.MethodPost, "/h/H1/device/lamp-1/toggle", nil)
	req.Header.Set("Datastar-Request", "true")
	req = req.WithContext(context.WithValue(req.Context(), authContextKey{}, AuthContext{UserID: user.ID, Username: user.Username}))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	body := w.Body.String()
	if !strings.Contains(body, `id="device-card-H1-lamp-1"`) || !strings.Contains(body, `data-toggle-error="Failed to toggle device."`) {
		t.Fatalf("expected card with datastar error marker, got %s", body)
	}
	if got := w.Header().Get("HX-Trigger"); got != "" {
		t.Fatalf("did not expect HX trigger, got %q", got)
	}
}

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
	"hata/internal/integrations"
	"hata/internal/util"

	"github.com/go-chi/chi/v5"
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
	if _, err := repos.Device().Create(context.Background(), "H1", "lamp-2", "Desk Lamp", "dummy", nil, "off"); err != nil {
		t.Fatalf("create offline device: %v", err)
	}
	if err := repos.Device().UpdateAvailability(context.Background(), "H1", "lamp-2", "offline"); err != nil {
		t.Fatalf("mark device offline: %v", err)
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
	if !strings.Contains(body, `href="/h/H1/manage"`) || !strings.Contains(body, `href="/h/H2/manage"`) {
		t.Fatalf("expected house manage links on page")
	}
	if !strings.Contains(body, "Bedroom Lamp") || !strings.Contains(body, "Desk Lamp") {
		t.Fatalf("expected house devices on page")
	}
	if !strings.Contains(body, "offline") {
		t.Fatalf("expected offline device status on page")
	}
	if !strings.Contains(body, "No devices in this house") {
		t.Fatalf("expected empty house message")
	}
}

func TestToggleHouseDevice_HTMXUpdatesOnlyDeviceCard(t *testing.T) {
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
	req.Header.Set("HX-Request", "true")
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
	if strings.Contains(body, "My devices") || strings.Contains(body, "Main Home") {
		t.Fatalf("expected only card partial, got full page: %s", body)
	}
	if strings.Contains(w.Header().Get("HX-Redirect"), "/app") {
		t.Fatalf("did not expect HX redirect")
	}
}

func TestHouseManagePage_RendersForMember(t *testing.T) {
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
	if _, err := repos.HouseDiscoveryNetwork().Create(context.Background(), "H1", "192.168.1.0/24", "Main WiFi"); err != nil {
		t.Fatalf("create discovery network: %v", err)
	}

	h := NewHandler(api.NewAuthHandler(repos), repos)
	router := chi.NewRouter()
	router.Get("/h/{houseId}/manage", h.HouseManagePage)

	req := httptest.NewRequest(http.MethodGet, "/h/H1/manage", nil)
	req = req.WithContext(context.WithValue(req.Context(), authContextKey{}, AuthContext{UserID: user.ID, Username: user.Username}))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "Main Home") || !strings.Contains(body, "House ID: H1") || !strings.Contains(body, "Your role: owner") {
		t.Fatalf("expected house manage details on page")
	}
	if !strings.Contains(body, "Devices") || !strings.Contains(body, "+ Add") || !strings.Contains(body, "Scan") {
		t.Fatalf("expected device discovery controls on page")
	}
	if !strings.Contains(body, "Discovery networks") || !strings.Contains(body, "Main WiFi") || !strings.Contains(body, "192.168.1.0/24") {
		t.Fatalf("expected discovery networks on page")
	}
}

func TestHouseDeviceDiscovery_StreamsFoundDevices(t *testing.T) {
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
	if _, err := repos.HouseDiscoveryNetwork().Create(context.Background(), "H1", "192.168.1.0/24", "Main WiFi"); err != nil {
		t.Fatalf("create discovery network: %v", err)
	}

	h := NewHandler(api.NewAuthHandler(repos), repos)
	var gotCIDRs []string
	h.discoverDevices = func(ctx context.Context, cidrs []string) <-chan integrations.DiscoveredDevice {
		gotCIDRs = append([]string(nil), cidrs...)
		ch := make(chan integrations.DiscoveredDevice, 1)
		ch <- integrations.DiscoveredDevice{Integration: "WiZ", Name: "Kitchen", IP: "192.168.1.41", State: "on"}
		close(ch)
		return ch
	}
	router := chi.NewRouter()
	router.Get("/h/{houseId}/manage/devices/discover", h.HouseDeviceDiscovery)

	req := httptest.NewRequest(http.MethodGet, "/h/H1/manage/devices/discover", nil)
	req = req.WithContext(context.WithValue(req.Context(), authContextKey{}, AuthContext{UserID: user.ID, Username: user.Username}))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if contentType := w.Header().Get("Content-Type"); !strings.Contains(contentType, "text/event-stream") {
		t.Fatalf("expected event stream content type, got %q", contentType)
	}
	body := w.Body.String()
	if !strings.Contains(body, "event: progress") || !strings.Contains(body, "Main WiFi: 192.168.1.0/24") || !strings.Contains(body, "event: device") || !strings.Contains(body, `"integration":"WiZ"`) || !strings.Contains(body, `"name":"Kitchen"`) || !strings.Contains(body, `"ip":"192.168.1.41"`) || !strings.Contains(body, `"state":"on"`) || !strings.Contains(body, `"inHouse":false`) {
		t.Fatalf("expected streamed discovery events, got %q", body)
	}
	if len(gotCIDRs) != 1 || gotCIDRs[0] != "192.168.1.0/24" {
		t.Fatalf("expected configured CIDR to be passed to discovery, got %v", gotCIDRs)
	}
}

func TestAddHouseDevice_AddsDiscoveredDevice(t *testing.T) {
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

	h := NewHandler(api.NewAuthHandler(repos), repos)
	router := chi.NewRouter()
	router.Post("/h/{houseId}/manage/devices", h.AddHouseDevice)

	req := httptest.NewRequest(http.MethodPost, "/h/H1/manage/devices", strings.NewReader("integration=WiZ&name=Kitchen&ip=192.168.1.41&state=off"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = req.WithContext(context.WithValue(req.Context(), authContextKey{}, AuthContext{UserID: user.ID, Username: user.Username}))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusSeeOther {
		t.Fatalf("expected 303, got %d: %s", w.Code, w.Body.String())
	}
	devices, err := repos.Device().ListByHouse(context.Background(), "H1")
	if err != nil {
		t.Fatalf("list devices: %v", err)
	}
	if len(devices) != 1 {
		t.Fatalf("expected 1 device, got %d", len(devices))
	}
	if devices[0].ID != "wiz-192-168-1-41" || devices[0].IntegrationID != "wiz" || devices[0].State != "off" || devices[0].IntegrationData == nil || !strings.Contains(*devices[0].IntegrationData, "192.168.1.41") {
		t.Fatalf("unexpected device: %#v", devices[0])
	}
}

func TestRenameHouseDevice_RenamesDeviceForManager(t *testing.T) {
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
	if _, err := repos.HouseRole().Assign(context.Background(), "H1", user.ID, "admin"); err != nil {
		t.Fatalf("assign house role: %v", err)
	}
	if _, err := repos.Device().Create(context.Background(), "H1", "lamp-1", "Old Lamp", "wiz", nil, "off"); err != nil {
		t.Fatalf("create device: %v", err)
	}

	h := NewHandler(api.NewAuthHandler(repos), repos)
	router := chi.NewRouter()
	router.Post("/h/{houseId}/manage/devices/{deviceId}/rename", h.RenameHouseDevice)

	req := httptest.NewRequest(http.MethodPost, "/h/H1/manage/devices/lamp-1/rename", strings.NewReader("name=New+Lamp"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = req.WithContext(context.WithValue(req.Context(), authContextKey{}, AuthContext{UserID: user.ID, Username: user.Username}))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusSeeOther {
		t.Fatalf("expected 303, got %d: %s", w.Code, w.Body.String())
	}
	device, err := repos.Device().GetByHouseAndID(context.Background(), "H1", "lamp-1")
	if err != nil {
		t.Fatalf("get device: %v", err)
	}
	if device == nil || device.Name != "New Lamp" {
		t.Fatalf("expected renamed device, got %#v", device)
	}
}

type fakeWebDeviceController struct{}

func (f *fakeWebDeviceController) SetState(ctx context.Context, device *db.DeviceData, state string) error {
	return nil
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

package webauth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"hata/internal/api"
	"hata/internal/integrations"
	"hata/internal/util"

	"github.com/go-chi/chi/v5"
)

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
	integrationData := `{"ip":"192.168.1.41"}`
	if _, err := repos.Device().Create(context.Background(), "H1", "lamp-1", "Office Lamp", "wiz", &integrationData, "off"); err != nil {
		t.Fatalf("create device: %v", err)
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
	if !strings.Contains(body, "House Settings") || !strings.Contains(body, "owner") || strings.Contains(body, "House ID:") || strings.Contains(body, "Your role:") {
		t.Fatalf("expected refreshed house manage details on page")
	}
	if !strings.Contains(body, "Office Lamp") || !strings.Contains(body, "IP 192.168.1.41") || !strings.Contains(body, "Rename Office Lamp") || !strings.Contains(body, "Remove Office Lamp") {
		t.Fatalf("expected refreshed device details on page")
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

func TestHouseDeviceDiscovery_SuggestsUpdateWhenMACMatchesChangedIP(t *testing.T) {
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
	integrationData := `{"ip":"192.168.1.41","mac":"aabbccddeeff"}`
	if _, err := repos.Device().Create(context.Background(), "H1", "office-light", "Office Light", "wiz", &integrationData, "off"); err != nil {
		t.Fatalf("create device: %v", err)
	}

	h := NewHandler(api.NewAuthHandler(repos), repos)
	h.discoverDevices = func(ctx context.Context, cidrs []string) <-chan integrations.DiscoveredDevice {
		ch := make(chan integrations.DiscoveredDevice, 1)
		ch <- integrations.DiscoveredDevice{Integration: "WiZ", Name: "Office Light", IP: "192.168.1.99", MAC: "aa:bb:cc:dd:ee:ff", State: "on"}
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
	body := w.Body.String()
	if !strings.Contains(body, `"inHouse":true`) || !strings.Contains(body, `"ipChanged":true`) || !strings.Contains(body, `"existingIp":"192.168.1.41"`) || !strings.Contains(body, `/h/H1/manage/devices/office-light/integration`) {
		t.Fatalf("expected streamed update suggestion, got %q", body)
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

	req := httptest.NewRequest(http.MethodPost, "/h/H1/manage/devices", strings.NewReader("integration=WiZ&name=Kitchen&ip=192.168.1.41&mac=aa%3Abb%3Acc%3Add%3Aee%3Aff&state=off"))
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
	if devices[0].ID != "wiz-192-168-1-41" || devices[0].IntegrationID != "wiz" || devices[0].State != "off" || devices[0].IntegrationData == nil || !strings.Contains(*devices[0].IntegrationData, "192.168.1.41") || !strings.Contains(*devices[0].IntegrationData, "aabbccddeeff") {
		t.Fatalf("unexpected device: %#v", devices[0])
	}
}

func TestUpdateHouseDeviceIntegration_UpdatesChangedIP(t *testing.T) {
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
	integrationData := `{"ip":"192.168.1.41","mac":"aabbccddeeff"}`
	if _, err := repos.Device().Create(context.Background(), "H1", "office-light", "Office Light", "wiz", &integrationData, "off"); err != nil {
		t.Fatalf("create device: %v", err)
	}

	h := NewHandler(api.NewAuthHandler(repos), repos)
	router := chi.NewRouter()
	router.Post("/h/{houseId}/manage/devices/{deviceId}/integration", h.UpdateHouseDeviceIntegration)

	req := httptest.NewRequest(http.MethodPost, "/h/H1/manage/devices/office-light/integration", strings.NewReader("integration=WiZ&ip=192.168.1.99&mac=aa%3Abb%3Acc%3Add%3Aee%3Aff&state=on"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = req.WithContext(context.WithValue(req.Context(), authContextKey{}, AuthContext{UserID: user.ID, Username: user.Username}))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusSeeOther {
		t.Fatalf("expected 303, got %d: %s", w.Code, w.Body.String())
	}
	device, err := repos.Device().GetByHouseAndID(context.Background(), "H1", "office-light")
	if err != nil {
		t.Fatalf("get device: %v", err)
	}
	if device == nil || device.IntegrationData == nil || !strings.Contains(*device.IntegrationData, "192.168.1.99") || !strings.Contains(*device.IntegrationData, "aabbccddeeff") || device.State != "on" || device.Availability != "online" {
		t.Fatalf("unexpected device after update: %#v", device)
	}
}

func TestDeleteHouseDevice_RemovesDeviceForManager(t *testing.T) {
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
	if _, err := repos.Device().Create(context.Background(), "H1", "lamp-1", "Lamp", "wiz", nil, "off"); err != nil {
		t.Fatalf("create device: %v", err)
	}

	h := NewHandler(api.NewAuthHandler(repos), repos)
	router := chi.NewRouter()
	router.Post("/h/{houseId}/manage/devices/{deviceId}/delete", h.DeleteHouseDevice)

	req := httptest.NewRequest(http.MethodPost, "/h/H1/manage/devices/lamp-1/delete", nil)
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
	if device != nil {
		t.Fatalf("expected deleted device, got %#v", device)
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

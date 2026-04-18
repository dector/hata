package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"hata/internal/db"
	"hata/internal/util"

	"github.com/go-chi/chi/v5"
)

func TestDeviceListByHouse_Unauthorized(t *testing.T) {
	repos, cleanup := setupAuthTest(t)
	defer cleanup()

	handler := NewDeviceHandler(repos)

	router := chi.NewRouter()
	router.Get("/api/latest/house/{houseId}/device", handler.ListByHouse)

	req := httptest.NewRequest(http.MethodGet, "/api/latest/house/house-1/device", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestDeviceListByHouse_Forbidden(t *testing.T) {
	repos, cleanup := setupAuthTest(t)
	defer cleanup()

	ctx := context.Background()
	user := createTestUser(t, repos, "user@example.com")
	token := createTestSession(t, repos, user.ID, strings.Repeat("a", 40))

	_, err := repos.House().Create(ctx, "house-1", "Main House")
	if err != nil {
		t.Fatalf("Failed to create house: %v", err)
	}

	handler := NewDeviceHandler(repos)

	router := chi.NewRouter()
	router.Get("/api/latest/house/{houseId}/device", handler.ListByHouse)

	req := httptest.NewRequest(http.MethodGet, "/api/latest/house/house-1/device", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}
}

func TestDeviceListByHouse_OK(t *testing.T) {
	repos, cleanup := setupAuthTest(t)
	defer cleanup()

	ctx := context.Background()
	user := createTestUser(t, repos, "user@example.com")
	token := createTestSession(t, repos, user.ID, strings.Repeat("b", 40))

	house, err := repos.House().Create(ctx, "house-1", "Main House")
	if err != nil {
		t.Fatalf("Failed to create house: %v", err)
	}

	_, err = repos.HouseRole().Assign(ctx, house.ID, user.ID, "owner")
	if err != nil {
		t.Fatalf("Failed to assign house role: %v", err)
	}

	integrationData := `{"type":"wifi","ip":"192.168.1.41"}`
	_, err = repos.Device().Create(ctx, house.ID, "lamp-office-1", "Office Lamp", "wiz", &integrationData, "on")
	if err != nil {
		t.Fatalf("Failed to create device: %v", err)
	}

	handler := NewDeviceHandler(repos)

	router := chi.NewRouter()
	router.Get("/api/latest/house/{houseId}/device", handler.ListByHouse)

	req := httptest.NewRequest(http.MethodGet, "/api/latest/house/house-1/device", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}

	var resp DeviceListByHouseResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(resp.Devices) != 1 {
		t.Fatalf("Expected 1 device, got %d", len(resp.Devices))
	}

	device := resp.Devices[0]
	if device.ID != "lamp-office-1" {
		t.Errorf("Expected device id 'lamp-office-1', got %q", device.ID)
	}
	if device.Name != "Office Lamp" {
		t.Errorf("Expected device name 'Office Lamp', got %q", device.Name)
	}
	if device.State != "on" {
		t.Errorf("Expected device state 'on', got %q", device.State)
	}
	if device.Integration.ID != "wiz" {
		t.Errorf("Expected integration id 'wiz', got %q", device.Integration.ID)
	}
	if device.Integration.Data == nil {
		t.Fatalf("Expected integration data, got nil")
	}
	if device.Integration.Data["type"] != "wifi" {
		t.Errorf("Expected integration type 'wifi', got %v", device.Integration.Data["type"])
	}
	if device.Integration.Data["ip"] != "192.168.1.41" {
		t.Errorf("Expected integration ip '192.168.1.41', got %v", device.Integration.Data["ip"])
	}
}

func TestDeviceListByUser_Unauthorized(t *testing.T) {
	repos, cleanup := setupAuthTest(t)
	defer cleanup()

	handler := NewDeviceHandler(repos)

	req := httptest.NewRequest(http.MethodGet, "/api/latest/device", nil)
	w := httptest.NewRecorder()

	handler.ListByUser(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestDeviceListByUser_OK(t *testing.T) {
	repos, cleanup := setupAuthTest(t)
	defer cleanup()

	ctx := context.Background()
	user := createTestUser(t, repos, "user@example.com")
	token := createTestSession(t, repos, user.ID, strings.Repeat("c", 40))

	houseA, err := repos.House().Create(ctx, "house-a", "House A")
	if err != nil {
		t.Fatalf("Failed to create house A: %v", err)
	}
	houseB, err := repos.House().Create(ctx, "house-b", "House B")
	if err != nil {
		t.Fatalf("Failed to create house B: %v", err)
	}

	_, err = repos.HouseRole().Assign(ctx, houseA.ID, user.ID, "owner")
	if err != nil {
		t.Fatalf("Failed to assign house role: %v", err)
	}

	validIntegration := `{"type":"wifi","ip":"192.168.1.99"}`
	_, err = repos.Device().Create(ctx, houseA.ID, "lamp-1", "Lamp", "wiz", &validIntegration, "off")
	if err != nil {
		t.Fatalf("Failed to create device: %v", err)
	}

	invalidIntegration := "not-json"
	_, err = repos.Device().Create(ctx, houseA.ID, "sensor-1", "Sensor", "wiz", &invalidIntegration, "")
	if err != nil {
		t.Fatalf("Failed to create device: %v", err)
	}

	_, err = repos.Device().Create(ctx, houseB.ID, "lamp-2", "Lamp 2", "wiz", &validIntegration, "on")
	if err != nil {
		t.Fatalf("Failed to create device: %v", err)
	}

	handler := NewDeviceHandler(repos)

	req := httptest.NewRequest(http.MethodGet, "/api/latest/device", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler.ListByUser(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}

	var resp DeviceListResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(resp.Devices) != 2 {
		t.Fatalf("Expected 2 devices, got %d", len(resp.Devices))
	}

	first := resp.Devices[0]
	second := resp.Devices[1]

	if first.House.ID != houseA.ID || second.House.ID != houseA.ID {
		t.Fatalf("Expected devices to belong to house %q", houseA.ID)
	}

	if first.ID != "lamp-1" || second.ID != "sensor-1" {
		t.Fatalf("Unexpected device ids: %q and %q", first.ID, second.ID)
	}

	if first.Integration.Data == nil || first.Integration.Data["ip"] != "192.168.1.99" {
		t.Errorf("Expected valid integration data for lamp-1")
	}
	if second.Integration.Data != nil {
		t.Errorf("Expected nil integration data for sensor-1")
	}
}

func TestDeviceSetState_Unauthorized(t *testing.T) {
	repos, cleanup := setupAuthTest(t)
	defer cleanup()

	handler := NewDeviceHandler(repos)
	router := chi.NewRouter()
	router.Patch("/api/latest/house/{houseId}/device/{deviceId}/state", handler.SetState)

	req := httptest.NewRequest(http.MethodPatch, "/api/latest/house/house-1/device/lamp-1/state", strings.NewReader(`{"state":"on"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("Expected status 401, got %d", w.Code)
	}
}

func TestDeviceSetState_Forbidden(t *testing.T) {
	repos, cleanup := setupAuthTest(t)
	defer cleanup()

	ctx := context.Background()
	user := createTestUser(t, repos, "user@example.com")
	token := createTestSession(t, repos, user.ID, strings.Repeat("d", 40))

	house, err := repos.House().Create(ctx, "house-1", "Main House")
	if err != nil {
		t.Fatalf("Failed to create house: %v", err)
	}
	_, err = repos.Device().Create(ctx, house.ID, "lamp-1", "Lamp", "wiz", nil, "off")
	if err != nil {
		t.Fatalf("Failed to create device: %v", err)
	}

	handler := NewDeviceHandler(repos)
	router := chi.NewRouter()
	router.Patch("/api/latest/house/{houseId}/device/{deviceId}/state", handler.SetState)

	req := httptest.NewRequest(http.MethodPatch, "/api/latest/house/house-1/device/lamp-1/state", strings.NewReader(`{"state":"on"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("Expected status 403, got %d", w.Code)
	}
}

func TestDeviceSetState_InvalidRequest(t *testing.T) {
	repos, cleanup := setupAuthTest(t)
	defer cleanup()

	ctx := context.Background()
	user := createTestUser(t, repos, "user@example.com")
	token := createTestSession(t, repos, user.ID, strings.Repeat("e", 40))

	house, err := repos.House().Create(ctx, "house-1", "Main House")
	if err != nil {
		t.Fatalf("Failed to create house: %v", err)
	}
	_, err = repos.HouseRole().Assign(ctx, house.ID, user.ID, "owner")
	if err != nil {
		t.Fatalf("Failed to assign role: %v", err)
	}
	_, err = repos.Device().Create(ctx, house.ID, "lamp-1", "Lamp", "wiz", nil, "off")
	if err != nil {
		t.Fatalf("Failed to create device: %v", err)
	}

	handler := NewDeviceHandler(repos)
	router := chi.NewRouter()
	router.Patch("/api/latest/house/{houseId}/device/{deviceId}/state", handler.SetState)

	t.Run("invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/api/latest/house/house-1/device/lamp-1/state", strings.NewReader(`{`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("invalid state", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/api/latest/house/house-1/device/lamp-1/state", strings.NewReader(`{"state":"maybe"}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Expected status 400, got %d", w.Code)
		}
	})
}

func TestDeviceSetState_NotFound(t *testing.T) {
	repos, cleanup := setupAuthTest(t)
	defer cleanup()

	ctx := context.Background()
	user := createTestUser(t, repos, "user@example.com")
	token := createTestSession(t, repos, user.ID, strings.Repeat("f", 40))

	house, err := repos.House().Create(ctx, "house-1", "Main House")
	if err != nil {
		t.Fatalf("Failed to create house: %v", err)
	}
	_, err = repos.HouseRole().Assign(ctx, house.ID, user.ID, "owner")
	if err != nil {
		t.Fatalf("Failed to assign role: %v", err)
	}

	handler := NewDeviceHandler(repos)
	router := chi.NewRouter()
	router.Patch("/api/latest/house/{houseId}/device/{deviceId}/state", handler.SetState)

	req := httptest.NewRequest(http.MethodPatch, "/api/latest/house/house-1/device/missing/state", strings.NewReader(`{"state":"on"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("Expected status 404, got %d", w.Code)
	}
}

func TestDeviceSetState_OK(t *testing.T) {
	repos, cleanup := setupAuthTest(t)
	defer cleanup()

	ctx := context.Background()
	user := createTestUser(t, repos, "user@example.com")
	token := createTestSession(t, repos, user.ID, strings.Repeat("g", 40))

	house, err := repos.House().Create(ctx, "house-1", "Main House")
	if err != nil {
		t.Fatalf("Failed to create house: %v", err)
	}
	_, err = repos.HouseRole().Assign(ctx, house.ID, user.ID, "owner")
	if err != nil {
		t.Fatalf("Failed to assign role: %v", err)
	}
	integrationData := `{"ip":"192.168.1.41"}`
	_, err = repos.Device().Create(ctx, house.ID, "lamp-1", "Lamp", "wiz", &integrationData, "off")
	if err != nil {
		t.Fatalf("Failed to create device: %v", err)
	}

	controller := &fakeDeviceController{}
	handler := NewDeviceHandlerWithController(repos, controller)
	router := chi.NewRouter()
	router.Patch("/api/latest/house/{houseId}/device/{deviceId}/state", handler.SetState)

	req := httptest.NewRequest(http.MethodPatch, "/api/latest/house/house-1/device/lamp-1/state", strings.NewReader(`{"state":"on"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}

	var resp DeviceSetStateResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if resp.DeviceID != "lamp-1" {
		t.Fatalf("Expected device id lamp-1, got %q", resp.DeviceID)
	}
	if resp.House.ID != house.ID {
		t.Fatalf("Expected house id %q, got %q", house.ID, resp.House.ID)
	}
	if resp.State != "on" {
		t.Fatalf("Expected state on, got %q", resp.State)
	}
	if controller.calls != 1 {
		t.Fatalf("Expected controller to be called once, got %d", controller.calls)
	}
	if controller.lastDeviceID != "lamp-1" {
		t.Fatalf("Expected controller device lamp-1, got %q", controller.lastDeviceID)
	}
	if controller.lastState != "on" {
		t.Fatalf("Expected controller state on, got %q", controller.lastState)
	}

	devices, err := repos.Device().ListByHouse(ctx, house.ID)
	if err != nil {
		t.Fatalf("Failed listing devices: %v", err)
	}
	if len(devices) != 1 {
		t.Fatalf("Expected 1 device, got %d", len(devices))
	}
	if devices[0].State != "on" {
		t.Fatalf("Expected persisted state on, got %q", devices[0].State)
	}
}

func TestDeviceSetState_ControlFailure_DoesNotPersist(t *testing.T) {
	repos, cleanup := setupAuthTest(t)
	defer cleanup()

	ctx := context.Background()
	user := createTestUser(t, repos, "user@example.com")
	token := createTestSession(t, repos, user.ID, strings.Repeat("h", 40))

	house, err := repos.House().Create(ctx, "house-1", "Main House")
	if err != nil {
		t.Fatalf("Failed to create house: %v", err)
	}
	_, err = repos.HouseRole().Assign(ctx, house.ID, user.ID, "owner")
	if err != nil {
		t.Fatalf("Failed to assign role: %v", err)
	}
	integrationData := `{"ip":"192.168.1.41"}`
	_, err = repos.Device().Create(ctx, house.ID, "lamp-1", "Lamp", "wiz", &integrationData, "off")
	if err != nil {
		t.Fatalf("Failed to create device: %v", err)
	}

	controller := &fakeDeviceController{err: errors.New("udp timeout")}
	handler := NewDeviceHandlerWithController(repos, controller)
	router := chi.NewRouter()
	router.Patch("/api/latest/house/{houseId}/device/{deviceId}/state", handler.SetState)

	req := httptest.NewRequest(http.MethodPatch, "/api/latest/house/house-1/device/lamp-1/state", strings.NewReader(`{"state":"on"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadGateway {
		t.Fatalf("Expected status 502, got %d", w.Code)
	}

	devices, err := repos.Device().ListByHouse(ctx, house.ID)
	if err != nil {
		t.Fatalf("Failed listing devices: %v", err)
	}
	if len(devices) != 1 {
		t.Fatalf("Expected 1 device, got %d", len(devices))
	}
	if devices[0].State != "off" {
		t.Fatalf("Expected persisted state to remain off, got %q", devices[0].State)
	}
}

type fakeDeviceController struct {
	calls        int
	lastDeviceID string
	lastState    string
	err          error
}

func (f *fakeDeviceController) SetState(ctx context.Context, device *db.DeviceData, state string) error {
	f.calls++
	f.lastDeviceID = device.ID
	f.lastState = state
	return f.err
}

func createTestUser(t *testing.T, repos db.Repositories, username string) *db.UserData {
	t.Helper()
	ctx := context.Background()
	passwordHash, err := util.HashPassword("testpass")
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}
	user, err := repos.User().Create(ctx, username, passwordHash, "Test User")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	return user
}

func createTestSession(t *testing.T, repos db.Repositories, userID int, token string) string {
	t.Helper()
	ctx := context.Background()
	_, err := repos.Session().Create(ctx, userID, token, time.Now().Add(2*time.Hour))
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}
	return token
}

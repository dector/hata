package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

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

func TestDeviceSetState_NoAck_ReturnsGatewayTimeout(t *testing.T) {
	repos, cleanup := setupAuthTest(t)
	defer cleanup()

	ctx := context.Background()
	user := createTestUser(t, repos, "user@example.com")
	token := createTestSession(t, repos, user.ID, strings.Repeat("i", 40))

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

	controller := &fakeDeviceController{err: errDeviceNoAck}
	handler := NewDeviceHandlerWithController(repos, controller)
	router := chi.NewRouter()
	router.Patch("/api/latest/house/{houseId}/device/{deviceId}/state", handler.SetState)

	req := httptest.NewRequest(http.MethodPatch, "/api/latest/house/house-1/device/lamp-1/state", strings.NewReader(`{"state":"on"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusGatewayTimeout {
		t.Fatalf("Expected status 504, got %d", w.Code)
	}

	var resp ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode error response: %v", err)
	}
	if resp.Error.Code != "device-no-ack" {
		t.Fatalf("Expected error code device-no-ack, got %q", resp.Error.Code)
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

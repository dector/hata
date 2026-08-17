package api

import (
	"context"
	"testing"
	"time"

	"hata/internal/db"
	"hata/internal/util"
)

type fakeDeviceController struct {
	calls           int
	lastDeviceID    string
	lastState       string
	lastBrightness  *int
	lastColorPreset *string
	err             error
}

func (f *fakeDeviceController) SetState(ctx context.Context, device *db.DeviceData, state string) error {
	f.calls++
	f.lastDeviceID = device.ID
	f.lastState = state
	return f.err
}

func (f *fakeDeviceController) SetLight(ctx context.Context, device *db.DeviceData, brightness *int, colorPreset *string) error {
	f.calls++
	f.lastDeviceID = device.ID
	f.lastBrightness = brightness
	f.lastColorPreset = colorPreset
	return f.err
}

func (f *fakeDeviceController) GetStatus(ctx context.Context, device *db.DeviceData) (DeviceStatus, error) {
	return DeviceStatus{State: f.lastState, Availability: "online"}, f.err
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

package api

import (
	"context"
	"errors"
	"testing"

	"hata/internal/db"
)

func TestRealDeviceController_SetState_UnsupportedIntegration(t *testing.T) {
	controller := NewRealDeviceController()
	device := &db.DeviceData{
		ID:            "plug-1",
		IntegrationID: "zigbee",
	}

	err := controller.SetState(context.Background(), device, "on")
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
	if !errors.Is(err, errUnsupportedIntegration) {
		t.Fatalf("Expected unsupported integration error, got %v", err)
	}
}

func TestRealDeviceController_SetState_WizMissingIP(t *testing.T) {
	controller := NewRealDeviceController()
	integrationData := `{"type":"wifi"}`
	device := &db.DeviceData{
		ID:              "lamp-1",
		IntegrationID:   "wiz",
		IntegrationData: &integrationData,
	}

	err := controller.SetState(context.Background(), device, "on")
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
	if !errors.Is(err, errInvalidIntegrationData) {
		t.Fatalf("Expected invalid integration data error, got %v", err)
	}
}

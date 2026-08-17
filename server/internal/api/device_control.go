package api

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"hata/internal/db"
)

var (
	errUnsupportedIntegration = errors.New("unsupported integration")
	errInvalidIntegrationData = errors.New("invalid integration data")
	errDeviceNoAck            = errors.New("device did not acknowledge command")
)

func IsDeviceNoAck(err error) bool {
	return errors.Is(err, errDeviceNoAck)
}

// DeviceStatus contains current physical device status.
type DeviceStatus struct {
	State            string
	Availability     string
	LightBrightness  *int
	LightColorPreset *string
}

// DeviceController executes physical control commands for devices.
type DeviceController interface {
	SetState(ctx context.Context, device *db.DeviceData, state string) error
	SetLight(ctx context.Context, device *db.DeviceData, brightness *int, colorPreset *string) error
	GetStatus(ctx context.Context, device *db.DeviceData) (DeviceStatus, error)
}

type RealDeviceController struct{}

func NewRealDeviceController() *RealDeviceController {
	return &RealDeviceController{}
}

func (c *RealDeviceController) SetState(ctx context.Context, device *db.DeviceData, state string) error {
	_ = ctx

	integrationID := strings.ToLower(strings.TrimSpace(device.IntegrationID))
	switch {
	case strings.HasPrefix(integrationID, "wiz"):
		return c.setWizState(device, state)
	default:
		return fmt.Errorf("%w: %s", errUnsupportedIntegration, device.IntegrationID)
	}
}

func (c *RealDeviceController) SetLight(ctx context.Context, device *db.DeviceData, brightness *int, colorPreset *string) error {
	integrationID := strings.ToLower(strings.TrimSpace(device.IntegrationID))
	switch {
	case strings.HasPrefix(integrationID, "wiz"):
		return c.setWizLight(ctx, device, brightness, colorPreset)
	default:
		return fmt.Errorf("%w: %s", errUnsupportedIntegration, device.IntegrationID)
	}
}

func (c *RealDeviceController) GetStatus(ctx context.Context, device *db.DeviceData) (DeviceStatus, error) {
	integrationID := strings.ToLower(strings.TrimSpace(device.IntegrationID))
	switch {
	case strings.HasPrefix(integrationID, "wiz"):
		return c.getWizStatus(ctx, device)
	default:
		return DeviceStatus{Availability: "unknown"}, fmt.Errorf("%w: %s", errUnsupportedIntegration, device.IntegrationID)
	}
}

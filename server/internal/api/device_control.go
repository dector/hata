package api

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"hata/internal/db"

	"github.com/dector/go-devices/pkg/wiz"
)

var (
	errUnsupportedIntegration = errors.New("unsupported integration")
	errInvalidIntegrationData = errors.New("invalid integration data")
)

// DeviceController executes physical control commands for devices.
type DeviceController interface {
	SetState(ctx context.Context, device *db.DeviceData, state string) error
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

func (c *RealDeviceController) setWizState(device *db.DeviceData, state string) error {
	integration := decodeIntegrationData(device.IntegrationData)
	if integration == nil {
		return fmt.Errorf("%w: missing integration data", errInvalidIntegrationData)
	}

	ip, _ := integration["ip"].(string)
	ip = strings.TrimSpace(ip)
	if ip == "" {
		return fmt.Errorf("%w: missing wiz ip", errInvalidIntegrationData)
	}

	light := wiz.Light(ip)
	switch state {
	case "on":
		if err := light.On(); err != nil {
			return fmt.Errorf("wiz turn on failed: %w", err)
		}
	case "off":
		if err := light.Off(); err != nil {
			return fmt.Errorf("wiz turn off failed: %w", err)
		}
	default:
		return fmt.Errorf("invalid state %q", state)
	}

	return nil
}

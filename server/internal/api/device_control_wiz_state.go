package api

import (
	"fmt"
	"strings"

	"hata/internal/db"

	"github.com/dector/go-devices/pkg/wiz"
)

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
			if wiz.IsNoAck(err) {
				return fmt.Errorf("%w: %w", errDeviceNoAck, err)
			}
			return fmt.Errorf("wiz turn on failed: %w", err)
		}
	case "off":
		if err := light.Off(); err != nil {
			if wiz.IsNoAck(err) {
				return fmt.Errorf("%w: %w", errDeviceNoAck, err)
			}
			return fmt.Errorf("wiz turn off failed: %w", err)
		}
	default:
		return fmt.Errorf("invalid state %q", state)
	}

	return nil
}

package api

import (
	"context"
	"fmt"
	"net"
	"strings"

	"hata/internal/db"
)

func (c *RealDeviceController) setWizLight(ctx context.Context, device *db.DeviceData, brightness *int, colorPreset *string) error {
	integration := decodeIntegrationData(device.IntegrationData)
	if integration == nil {
		return fmt.Errorf("%w: missing integration data", errInvalidIntegrationData)
	}

	ip, _ := integration["ip"].(string)
	ip = strings.TrimSpace(ip)
	if ip == "" {
		return fmt.Errorf("%w: missing wiz ip", errInvalidIntegrationData)
	}

	params := map[string]any{}
	if brightness != nil {
		params["dimming"] = *brightness
	}
	if colorPreset != nil {
		for key, value := range wizPresetParams(*colorPreset) {
			params[key] = value
		}
	}

	if _, err := udpJSON(ctx, net.JoinHostPort(ip, "38899"), map[string]any{"method": "setPilot", "params": params}); err != nil {
		return fmt.Errorf("%w: %w", errDeviceNoAck, err)
	}
	return nil
}

func wizPresetParams(preset string) map[string]any {
	switch preset {
	case "warm_white":
		return map[string]any{"temp": 2700}
	case "soft_white":
		return map[string]any{"temp": 3000}
	case "daylight_white":
		return map[string]any{"temp": 5000}
	case "cold_white":
		return map[string]any{"temp": 6500}
	case "red":
		return map[string]any{"r": 255, "g": 0, "b": 0}
	case "orange":
		return map[string]any{"r": 255, "g": 128, "b": 0}
	case "yellow":
		return map[string]any{"r": 255, "g": 220, "b": 0}
	case "green":
		return map[string]any{"r": 0, "g": 255, "b": 0}
	case "blue":
		return map[string]any{"r": 0, "g": 80, "b": 255}
	case "purple":
		return map[string]any{"r": 128, "g": 0, "b": 255}
	default:
		return map[string]any{}
	}
}

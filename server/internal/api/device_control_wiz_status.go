package api

import (
	"context"
	"fmt"
	"net"
	"strings"

	"hata/internal/db"
)

func (c *RealDeviceController) getWizStatus(ctx context.Context, device *db.DeviceData) (DeviceStatus, error) {
	integration := decodeIntegrationData(device.IntegrationData)
	if integration == nil {
		return DeviceStatus{Availability: "unknown"}, fmt.Errorf("%w: missing integration data", errInvalidIntegrationData)
	}

	ip, _ := integration["ip"].(string)
	ip = strings.TrimSpace(ip)
	if ip == "" {
		return DeviceStatus{Availability: "unknown"}, fmt.Errorf("%w: missing wiz ip", errInvalidIntegrationData)
	}

	response, err := udpJSON(ctx, net.JoinHostPort(ip, "38899"), map[string]any{"method": "getPilot", "params": map[string]any{}})
	if err != nil {
		return DeviceStatus{Availability: "offline"}, fmt.Errorf("%w: %w", errDeviceNoAck, err)
	}

	result, _ := response["result"].(map[string]any)
	state, ok := result["state"].(bool)
	if !ok {
		return DeviceStatus{Availability: "online"}, nil
	}
	status := DeviceStatus{Availability: "online"}
	if state {
		status.State = "on"
	} else {
		status.State = "off"
	}
	if dimming, ok := numberAsInt(result["dimming"]); ok {
		status.LightBrightness = &dimming
	}
	return status, nil
}
func numberAsInt(value any) (int, bool) {
	switch v := value.(type) {
	case float64:
		return int(v), true
	case int:
		return v, true
	default:
		return 0, false
	}
}

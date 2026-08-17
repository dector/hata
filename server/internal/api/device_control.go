package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"hata/internal/db"

	"github.com/dector/go-devices/pkg/wiz"
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

func udpJSON(ctx context.Context, address string, payload any) (map[string]any, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	dialer := net.Dialer{Timeout: time.Second}
	conn, err := dialer.DialContext(ctx, "udp", address)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	_ = conn.SetDeadline(time.Now().Add(time.Second))
	if _, err := conn.Write(body); err != nil {
		return nil, err
	}

	buf := make([]byte, 2048)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}

	var response map[string]any
	if err := json.Unmarshal(buf[:n], &response); err != nil {
		return nil, err
	}
	return response, nil
}

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
	State        string
	Availability string
}

// DeviceController executes physical control commands for devices.
type DeviceController interface {
	SetState(ctx context.Context, device *db.DeviceData, state string) error
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
	if state {
		return DeviceStatus{State: "on", Availability: "online"}, nil
	}
	return DeviceStatus{State: "off", Availability: "online"}, nil
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

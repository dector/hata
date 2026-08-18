package integrations

import (
	"context"
	"encoding/json"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/dector/go-devices/pkg/wiz"
)

const wizDiscoveryTimeout = 4 * time.Second

// DiscoveredDevice is a device found by an integration discovery scan.
type DiscoveredDevice struct {
	Integration string `json:"integration"`
	Name        string `json:"name"`
	IP          string `json:"ip"`
	MAC         string `json:"mac"`
	State       string `json:"state"`
}

// DiscoverDevices runs discovery for all supported integrations.
func DiscoverDevices(ctx context.Context, cidrs []string) <-chan DiscoveredDevice {
	out := make(chan DiscoveredDevice)
	go func() {
		defer close(out)

		seen := sync.Map{}

		for _, cidr := range uniqueCIDRs(cidrs) {
			discoverWiZ(ctx, wiz.DiscoverDevicesInSubnet(cidr, wizDiscoveryTimeout), out, &seen)
		}
		discoverWiZ(ctx, wiz.DiscoverDevices(wizDiscoveryTimeout), out, &seen)
	}()
	return out
}

func discoverWiZ(ctx context.Context, devices <-chan wiz.DiscoveredDevice, out chan<- DiscoveredDevice, seen *sync.Map) {
	for device := range devices {
		if _, loaded := seen.LoadOrStore("wiz:"+device.IP, struct{}{}); loaded {
			continue
		}

		name := strings.TrimSpace(device.ModuleName)
		if name == "" {
			name = strings.TrimSpace(device.MAC)
		}

		state := wizState(ctx, device.IP)

		select {
		case <-ctx.Done():
			return
		case out <- DiscoveredDevice{Integration: "WiZ", Name: name, IP: device.IP, MAC: normalizeMAC(device.MAC), State: state}:
		}
	}
}

func normalizeMAC(mac string) string {
	return strings.ToLower(strings.TrimSpace(mac))
}

func wizState(ctx context.Context, ip string) string {
	response, err := udpJSON(ctx, net.JoinHostPort(ip, "38899"), map[string]any{"method": "getPilot", "params": map[string]any{}})
	if err != nil {
		return ""
	}
	result, _ := response["result"].(map[string]any)
	state, ok := result["state"].(bool)
	if !ok {
		return ""
	}
	if state {
		return "on"
	}
	return "off"
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

func uniqueCIDRs(cidrs []string) []string {
	out := make([]string, 0, len(cidrs))
	seen := map[string]struct{}{}
	for _, cidr := range cidrs {
		cidr = strings.TrimSpace(cidr)
		if cidr == "" {
			continue
		}
		if _, ok := seen[cidr]; ok {
			continue
		}
		seen[cidr] = struct{}{}
		out = append(out, cidr)
	}
	return out
}

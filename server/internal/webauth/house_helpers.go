package webauth

import (
	"encoding/json"
	"fmt"
	"hata/internal/db"
	"net/http"
	"net/url"
	"strings"
)

func writeSSE(w http.ResponseWriter, event, data string) {
	fmt.Fprintf(w, "event: %s\n", event)
	for line := range strings.SplitSeq(data, "\n") {
		fmt.Fprintf(w, "data: %s\n", line)
	}
	fmt.Fprint(w, "\n")
}

func deviceIPsByIntegration(devices []*db.DeviceData) map[string]bool {
	result := map[string]bool{}
	for key := range devicesByIntegrationField(devices, "ip") {
		result[key] = true
	}
	return result
}

func devicesByIntegrationField(devices []*db.DeviceData, field string) map[string]*db.DeviceData {
	result := map[string]*db.DeviceData{}
	for _, device := range devices {
		value := strings.TrimSpace(deviceIntegrationField(device, field))
		if value == "" {
			continue
		}
		if field == "mac" {
			value = normalizeMAC(value)
		}
		result[strings.ToLower(device.IntegrationID)+":"+value] = device
	}
	return result
}

func deviceIntegrationField(device *db.DeviceData, field string) string {
	if device == nil || device.IntegrationData == nil {
		return ""
	}
	var data map[string]any
	if err := json.Unmarshal([]byte(*device.IntegrationData), &data); err != nil {
		return ""
	}
	value, _ := data[field].(string)
	return strings.TrimSpace(value)
}

func normalizeMAC(mac string) string {
	mac = strings.ToLower(strings.TrimSpace(mac))
	mac = strings.ReplaceAll(mac, ":", "")
	mac = strings.ReplaceAll(mac, "-", "")
	return mac
}

func deviceIDForIntegrationIP(integration, ip string) string {
	id := strings.ToLower(strings.TrimSpace(integration)) + "-" + strings.TrimSpace(ip)
	id = strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z':
			return r
		case r >= '0' && r <= '9':
			return r
		default:
			return '-'
		}
	}, id)
	id = strings.Trim(id, "-")
	if id == "" {
		return "device"
	}
	return id
}

func canManageHouse(role string) bool {
	switch strings.ToLower(strings.TrimSpace(role)) {
	case "owner", "admin":
		return true
	default:
		return false
	}
}

func houseManageURL(houseID string) string {
	return "/h/" + url.PathEscape(houseID) + "/manage"
}

func discoveryProgressMessage(networks []*db.HouseDiscoveryNetworkData) string {
	parts := []string{"Scanning WiZ devices on server local networks"}
	for _, network := range networks {
		label := strings.TrimSpace(network.Label)
		if label == "" {
			parts = append(parts, network.CIDR)
			continue
		}
		parts = append(parts, label+": "+network.CIDR)
	}
	return strings.Join(parts, "; ") + "…"
}

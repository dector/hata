package webui

import (
	"net/url"
	"strconv"
	"strings"
)

const LocalHTMXScriptPath = "/assets/js/htmx-2.0.4.min.js"

type HomePageData struct {
	IsLoggedIn    bool
	DisplayName   string
	HTMXScriptSrc string
}

type LoginPageData struct {
	Then          string
	Error         string
	HTMXScriptSrc string
}

type LogoutPageData struct {
	HTMXScriptSrc string
}

type MePageData struct {
	UserID          int
	Username        string
	DisplayName     string
	HasSessionToken bool
	HTMXScriptSrc   string
}

type AppPageData struct {
	DisplayName   string
	Houses        []AppHouseData
	HTMXScriptSrc string
}

type AppHouseData struct {
	ID          string
	DisplayName string
	Role        string
	Devices     []AppDeviceData
}

type AppDeviceData struct {
	HouseID          string
	ID               string
	Name             string
	IntegrationID    string
	State            string
	Availability     string
	ToggleURL        string
	StateURL         string
	LightURL         string
	IsLight          bool
	LightBrightness  *int
	LightColorPreset *string
	LightPresets     []LightPresetData
}

type LightPresetData struct {
	ID    string
	Label string
	Hex   string
}

type HouseManagePageData struct {
	DisplayName       string
	House             AppHouseData
	DiscoveryNetworks []HouseDiscoveryNetworkData
	CanManageHouse    bool
	HTMXScriptSrc     string
}

type HouseDiscoveryNetworkData struct {
	ID    int
	CIDR  string
	Label string
}

func lightDialogID(device AppDeviceData) string {
	id := device.HouseID + "-" + device.ID
	id = strings.NewReplacer("/", "-", " ", "-", ".", "-", ":", "-").Replace(id)
	return "light-" + id
}

func lightBrightnessValue(device AppDeviceData) string {
	if device.LightBrightness == nil {
		return "50"
	}
	return strconv.Itoa(*device.LightBrightness)
}

func lightPresetChecked(device AppDeviceData, presetID string) bool {
	return device.LightColorPreset != nil && *device.LightColorPreset == presetID
}

func deviceStatusClass(device AppDeviceData) string {
	if strings.EqualFold(device.Availability, "offline") {
		return "device-status-offline"
	}

	switch strings.ToLower(strings.TrimSpace(device.State)) {
	case "on":
		return "device-status-on"
	case "off":
		return "device-status-off"
	default:
		return "device-status-unknown"
	}
}

func deviceStatusLabel(device AppDeviceData) string {
	if strings.EqualFold(device.Availability, "offline") {
		return "offline"
	}

	switch strings.ToLower(strings.TrimSpace(device.State)) {
	case "on":
		return "on"
	case "off":
		return "off"
	default:
		return device.State
	}
}

func scriptSrcOrDefault(src string) string {
	if strings.TrimSpace(src) == "" {
		return LocalHTMXScriptPath
	}
	return src
}

func houseManagePath(id string) string {
	return "/h/" + url.PathEscape(id) + "/manage"
}

func houseManageDiscoverPath(id string) string {
	return houseManagePath(id) + "/devices/discover"
}

func houseDeviceRenamePath(houseID string, deviceID string) string {
	return houseManagePath(houseID) + "/devices/" + url.PathEscape(deviceID) + "/rename"
}

func HouseDeviceTogglePath(houseID string, deviceID string) string {
	return "/h/" + url.PathEscape(houseID) + "/device/" + url.PathEscape(deviceID) + "/toggle"
}

func HouseDeviceStatePath(houseID string, deviceID string) string {
	return "/h/" + url.PathEscape(houseID) + "/device/" + url.PathEscape(deviceID) + "/state"
}

func HouseDeviceLightPath(houseID string, deviceID string) string {
	return "/h/" + url.PathEscape(houseID) + "/device/" + url.PathEscape(deviceID) + "/light"
}

func houseDiscoveryNetworksPath(id string) string {
	return houseManagePath(id) + "/discovery-networks"
}

func houseDiscoveryNetworkDeletePath(houseID string, networkID int) string {
	return houseDiscoveryNetworksPath(houseID) + "/" + strconv.Itoa(networkID) + "/delete"
}

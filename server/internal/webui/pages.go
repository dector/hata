package webui

import (
	"net/url"
	"strconv"
	"strings"
)

const LocalDatastarScriptPath = "/assets/js/datastar-1.0.2.js"

type HomePageData struct {
	IsLoggedIn  bool
	DisplayName string
}

type LoginPageData struct {
	Then  string
	Error string
}

type LogoutPageData struct{}

type MePageData struct {
	UserID          int
	Username        string
	DisplayName     string
	HasSessionToken bool
}

type AppPageData struct {
	DisplayName   string
	Houses        []AppHouseData
	HeaderHouses  []AppHouseData
	ActiveHouseID string
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
	IntegrationIP    string
	State            string
	Availability     string
	ToggleURL        string
	StateURL         string
	LightURL         string
	IsLight          bool
	LightBrightness  *int
	LightColorPreset *string
	LightPresets     []LightPresetData
	ToggleError      string
}

type LightPresetData struct {
	ID    string
	Label string
	Hex   string
}

type HouseManagePageData struct {
	DisplayName       string
	House             AppHouseData
	HeaderHouses             []AppHouseData
	ActiveHouseID            string
	DefaultDiscoveryNetworks []string
	DiscoveryNetworks        []HouseDiscoveryNetworkData
	CanManageHouse    bool
}

type HouseDiscoveryNetworkData struct {
	ID    int
	CIDR  string
	Label string
}

func houseDevicesTitle(house AppHouseData) string {
	return "Devices (" + strconv.Itoa(len(house.Devices)) + ")"
}

func activeHouseName(houses []AppHouseData, activeHouseID string) string {
	for _, house := range houses {
		if house.ID == activeHouseID {
			return house.DisplayName
		}
	}
	if len(houses) > 0 {
		return houses[0].DisplayName
	}
	return "House"
}

func activeHouseRole(houses []AppHouseData, activeHouseID string) string {
	for _, house := range houses {
		if house.ID == activeHouseID {
			return house.Role
		}
	}
	if len(houses) > 0 {
		return houses[0].Role
	}
	return ""
}

func lightDialogID(device AppDeviceData) string {
	return "light-" + stableDeviceElementID(device)
}

func deviceCardID(device AppDeviceData) string {
	return "device-card-" + stableDeviceElementID(device)
}

func datastarPostFormAction(url string, selector string) string {
	return "@post(" + strconv.Quote(url) + ", {contentType: 'form', selector: " + strconv.Quote(selector) + "})"
}

func stableDeviceElementID(device AppDeviceData) string {
	id := device.HouseID + "-" + device.ID
	id = strings.NewReplacer("/", "-", " ", "-", ".", "-", ":", "-").Replace(id)
	return id
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

func toggleNextState(device AppDeviceData) string {
	if strings.EqualFold(strings.TrimSpace(device.State), "on") {
		return "off"
	}
	return "on"
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

func activeHousePath() string {
	return "/app/active-house"
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

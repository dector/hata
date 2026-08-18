package webui

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
	"time"
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
	ID            string
	DisplayName   string
	Role          string
	Weather       *AppWeatherData
	Devices       []AppDeviceData
	ShoppingLists []ShoppingListData
}

// AppWeatherData contains weather shown on the app page.
type AppWeatherData struct {
	Status          string
	LocationLabel   string
	Temperature     float64
	TemperatureUnit string
	ConditionText   string
	ConditionIcon   string
	HumidityPercent *float64
	WindSpeed       *float64
	WindSpeedUnit   string
	UpdatedAt       time.Time
}

type AppDeviceData struct {
	HouseID          string
	ID               string
	Name             string
	IntegrationID    string
	IntegrationIP    string
	IntegrationData  string
	State            string
	Availability     string
	DetailsURL       string
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
	DisplayName              string
	House                    AppHouseData
	HeaderHouses             []AppHouseData
	ActiveHouseID            string
	DefaultDiscoveryNetworks []string
	DiscoveryNetworks        []HouseDiscoveryNetworkData
	CanManageHouse           bool
}

type DevicePageData struct {
	DisplayName   string
	HeaderHouses  []AppHouseData
	ActiveHouseID string
	House         AppHouseData
	Device        AppDeviceData
	ReloadURL     string
	ReloadMessage string
	ReloadError   string
}

type ShoppingListPageData struct {
	DisplayName    string
	HeaderHouses   []AppHouseData
	ActiveHouseID  string
	House          AppHouseData
	List           ShoppingListData
	Items          []ShoppingItemData
	PageURL        string
	CanManageHouse bool
}

type ShoppingListData struct {
	ID   string
	Name string
}

type ShoppingItemData struct {
	ID      string
	Name    string
	Checked bool
}

type HouseDiscoveryNetworkData struct {
	ID    int
	CIDR  string
	Label string
}

func houseDevicesTitle(house AppHouseData) string {
	return "Devices (" + strconv.Itoa(len(house.Devices)) + ")"
}

func hasVisibleWeather(house AppHouseData) bool {
	return house.Weather != nil && (house.Weather.Status == "ok" || house.Weather.Status == "stale")
}

func weatherIconPath(icon string) string {
	switch strings.TrimSpace(icon) {
	case "clear":
		return "M120,40V16a8,8,0,0,1,16,0V40a8,8,0,0,1-16,0Zm72,88a64,64,0,1,1-64-64A64.07,64.07,0,0,1,192,128Zm-16,0a48,48,0,1,0-48,48A48.05,48.05,0,0,0,176,128ZM58.34,69.66A8,8,0,0,0,69.66,58.34l-16-16A8,8,0,0,0,42.34,53.66Zm0,116.68-16,16a8,8,0,0,0,11.32,11.32l16-16a8,8,0,0,0-11.32-11.32ZM192,72a8,8,0,0,0,5.66-2.34l16-16a8,8,0,0,0-11.32-11.32l-16,16A8,8,0,0,0,192,72Zm5.66,114.34a8,8,0,0,0-11.32,11.32l16,16a8,8,0,0,0,11.32-11.32ZM48,128a8,8,0,0,0-8-8H16a8,8,0,0,0,0,16H40A8,8,0,0,0,48,128Zm80,80a8,8,0,0,0-8,8v24a8,8,0,0,0,16,0V216A8,8,0,0,0,128,208Zm112-88H216a8,8,0,0,0,0,16h24a8,8,0,0,0,0-16Z"
	case "mostly_clear", "partly_cloudy":
		return "M164,72a76.2,76.2,0,0,0-20.26,2.73,55.63,55.63,0,0,0-9.41-11.54l9.51-13.57a8,8,0,1,0-13.11-9.18L121.22,54A55.9,55.9,0,0,0,96,48c-.58,0-1.16,0-1.74,0L91.37,31.71a8,8,0,1,0-15.75,2.77L78.5,50.82A56.1,56.1,0,0,0,55.23,65.67L41.61,56.14a8,8,0,1,0-9.17,13.11L46,78.77A55.55,55.55,0,0,0,40,104c0,.57,0,1.15,0,1.72L23.71,108.6a8,8,0,0,0,1.38,15.88,8.24,8.24,0,0,0,1.39-.12l16.32-2.88a55.74,55.74,0,0,0,5.86,12.42A52,52,0,0,0,84,224h80a76,76,0,0,0,0-152ZM56,104a40,40,0,0,1,72.54-23.24,76.26,76.26,0,0,0-35.62,40,52.14,52.14,0,0,0-31,4.17A40,40,0,0,1,56,104ZM164,208H84a36,36,0,1,1,4.78-71.69c-.37,2.37-.63,4.79-.77,7.23a8,8,0,0,0,16,.92,58.91,58.91,0,0,1,1.88-11.81c0-.16.09-.32.12-.48A60.06,60.06,0,1,1,164,208Z"
	case "cloudy":
		return "M160,40A88.09,88.09,0,0,0,81.29,88.67,64,64,0,1,0,72,216h88a88,88,0,0,0,0-176Zm0,160H72a48,48,0,0,1,0-96c1.1,0,2.2,0,3.29.11A88,88,0,0,0,72,128a8,8,0,0,0,16,0,72,72,0,1,1,72,72Z"
	case "fog":
		return "M120,208H72a8,8,0,0,1,0-16h48a8,8,0,0,1,0,16Zm64-16H160a8,8,0,0,0,0,16h24a8,8,0,0,0,0-16Zm-24,32H104a8,8,0,0,0,0,16h56a8,8,0,0,0,0-16Zm72-124a76.08,76.08,0,0,1-76,76H76A52,52,0,0,1,76,72a53.26,53.26,0,0,1,8.92.76A76.08,76.08,0,0,1,232,100Zm-16,0A60.06,60.06,0,0,0,96,96.46a8,8,0,0,1-16-.92q.21-3.66.77-7.23A38.11,38.11,0,0,0,76,88a36,36,0,0,0,0,72h80A60.07,60.07,0,0,0,216,100Z"
	case "drizzle", "freezing_drizzle", "rain", "freezing_rain", "rain_showers":
		return "M158.66,196.44l-32,48a8,8,0,1,1-13.32-8.88l32-48a8,8,0,0,1,13.32,8.88ZM232,92a76.08,76.08,0,0,1-76,76H132.28l-29.62,44.44a8,8,0,1,1-13.32-8.88L113.05,168H76A52,52,0,0,1,76,64a53.26,53.26,0,0,1,8.92.76A76.08,76.08,0,0,1,232,92Zm-16,0A60.06,60.06,0,0,0,96,88.46a8,8,0,0,1-16-.92q.21-3.66.77-7.23A38.11,38.11,0,0,0,76,80a36,36,0,0,0,0,72h80A60.07,60.07,0,0,0,216,92Z"
	case "snow", "snow_showers":
		return "M88,196a12,12,0,1,1-12-12A12,12,0,0,1,88,196Zm28,4a12,12,0,1,0,12,12A12,12,0,0,0,116,200Zm48-16a12,12,0,1,0,12,12A12,12,0,0,0,164,184ZM68,224a12,12,0,1,0,12,12A12,12,0,0,0,68,224Zm88,0a12,12,0,1,0,12,12A12,12,0,0,0,156,224ZM232,92a76.08,76.08,0,0,1-76,76H76A52,52,0,0,1,76,64a53.26,53.26,0,0,1,8.92.76A76.08,76.08,0,0,1,232,92Zm-16,0A60.06,60.06,0,0,0,96,88.46a8,8,0,0,1-16-.92q.21-3.66.77-7.23A38.11,38.11,0,0,0,76,80a36,36,0,0,0,0,72h80A60.07,60.07,0,0,0,216,92Z"
	case "thunderstorm", "thunderstorm_hail":
		return "M156,16A76.2,76.2,0,0,0,84.92,64.76,53.26,53.26,0,0,0,76,64a52,52,0,0,0,0,104h37.87L97.14,195.88A8,8,0,0,0,104,208h25.87l-16.73,27.88a8,8,0,0,0,13.72,8.24l24-40A8,8,0,0,0,144,192H118.13l14.4-24H156a76,76,0,0,0,0-152Zm0,136H76a36,36,0,0,1,0-72,38.11,38.11,0,0,1,4.78.31q-.56,3.57-.77,7.23a8,8,0,0,0,16,.92A60.06,60.06,0,1,1,156,152Z"
	default:
		return "M140,180a12,12,0,1,1-12-12A12,12,0,0,1,140,180ZM128,72c-22.06,0-40,16.15-40,36v4a8,8,0,0,0,16,0v-4c0-11,10.77-20,24-20s24,9,24,20-10.77,20-24,20a8,8,0,0,0-8,8v8a8,8,0,0,0,16,0v-.72c18.24-3.35,32-17.9,32-35.28C168,88.15,150.06,72,128,72Zm104,56A104,104,0,1,1,128,24,104.11,104.11,0,0,1,232,128Zm-16,0a88,88,0,1,0-88,88A88.1,88.1,0,0,0,216,128Z"
	}
}

func weatherTemperatureLabel(weather AppWeatherData) string {
	unit := strings.TrimSpace(weather.TemperatureUnit)
	if unit == "" {
		unit = "°C"
	}
	if strings.HasPrefix(unit, "°") {
		return strconv.FormatFloat(weather.Temperature, 'f', 1, 64) + " " + unit
	}
	return strconv.FormatFloat(weather.Temperature, 'f', 1, 64) + " °" + unit
}

func weatherUpdatedLabel(weather AppWeatherData) string {
	if weather.UpdatedAt.IsZero() {
		return "updated unknown"
	}
	age := time.Since(weather.UpdatedAt)
	if age < 0 {
		age = 0
	}
	minutes := int(age.Minutes())
	if minutes < 1 {
		return "updated just now"
	}
	if minutes < 60 {
		if minutes == 1 {
			return "updated 1 minute ago"
		}
		return "updated " + strconv.Itoa(minutes) + " minutes ago"
	}
	hours := int(age.Hours())
	if hours < 24 {
		if hours == 1 {
			return "updated 1 hour ago"
		}
		return "updated " + strconv.Itoa(hours) + " hours ago"
	}
	days := hours / 24
	if days == 1 {
		return "updated 1 day ago"
	}
	return "updated " + strconv.Itoa(days) + " days ago"
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
		if strings.TrimSpace(device.State) == "" {
			return "unknown"
		}
		return device.State
	}
}

func formattedIntegrationData(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "No integration metadata stored."
	}
	var v any
	if err := json.Unmarshal([]byte(raw), &v); err != nil {
		return raw
	}
	buf, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return raw
	}
	return string(buf)
}

func nullableStringValue(value *string) string {
	if value == nil || strings.TrimSpace(*value) == "" {
		return "unknown"
	}
	return *value
}

func nullableIntValue(value *int) string {
	if value == nil {
		return "unknown"
	}
	return strconv.Itoa(*value)
}

func lightSupportLabel(device AppDeviceData) string {
	if device.IsLight {
		return "yes"
	}
	return "no"
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

func houseDeviceDeletePath(houseID string, deviceID string) string {
	return houseManagePath(houseID) + "/devices/" + url.PathEscape(deviceID) + "/delete"
}

func DeviceDetailsPath(deviceID string) string {
	return "/d/" + url.PathEscape(deviceID)
}

func DeviceReloadPath(deviceID string) string {
	return DeviceDetailsPath(deviceID) + "/reload"
}

func ShoppingListPath(listID string) string {
	return "/sl/" + url.PathEscape(listID)
}

func shoppingListItemsPath(listID string) string {
	return ShoppingListPath(listID) + "/items"
}

func shoppingListItemCheckPath(listID string, itemID string) string {
	return ShoppingListPath(listID) + "/items/" + url.PathEscape(itemID) + "/check"
}

func shoppingListItemRenamePath(listID string, itemID string) string {
	return ShoppingListPath(listID) + "/items/" + url.PathEscape(itemID) + "/rename"
}

func shoppingListItemDeletePath(listID string, itemID string) string {
	return ShoppingListPath(listID) + "/items/" + url.PathEscape(itemID) + "/delete"
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

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
	ID            string
	Name          string
	IntegrationID string
	State         string
	Availability  string
	ToggleURL     string
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

func houseDiscoveryNetworksPath(id string) string {
	return houseManagePath(id) + "/discovery-networks"
}

func houseDiscoveryNetworkDeletePath(houseID string, networkID int) string {
	return houseDiscoveryNetworksPath(houseID) + "/" + strconv.Itoa(networkID) + "/delete"
}

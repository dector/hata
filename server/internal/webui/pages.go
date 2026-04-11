package webui

import "strings"

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
	HasSessionToken bool
	HTMXScriptSrc   string
}

type AppPageData struct {
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
}

func scriptSrcOrDefault(src string) string {
	if strings.TrimSpace(src) == "" {
		return LocalHTMXScriptPath
	}
	return src
}

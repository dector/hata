package webui

import "strings"

const LocalHTMXScriptPath = "/assets/js/htmx-2.0.4.min.js"

type HomePageData struct {
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

func scriptSrcOrDefault(src string) string {
	if strings.TrimSpace(src) == "" {
		return LocalHTMXScriptPath
	}
	return src
}

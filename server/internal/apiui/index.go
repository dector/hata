package apiui

import (
	"html/template"
	"net/http"
)

// IndexHandler serves the API UI index page.
type IndexHandler struct{}

// NewIndexHandler creates a new IndexHandler.
func NewIndexHandler() *IndexHandler {
	return &IndexHandler{}
}

// Index serves a list of available /apiui endpoints.
func (h *IndexHandler) Index(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	writePage(w, pageViewData{
		Title:       "Hata API UI",
		Heading:     "Hata API UI",
		Description: template.HTML(`Available browser test pages under <code>/apiui</code>.`),
		Path:        "/apiui",
		BackHref:    "/dev",
		Content: template.HTML(`
<p style="margin-top:0">Available pages</p>
<ul>
  <li><a href="/apiui/latest/auth/login">/apiui/latest/auth/login</a> — <code>POST /api/latest/auth/login</code></li>
  <li><a href="/apiui/latest/house">/apiui/latest/house</a> — <code>GET /api/latest/house</code></li>
  <li><a href="/apiui/latest/device">/apiui/latest/device</a> — <code>GET /api/latest/device</code></li>
  <li><a href="/apiui/latest/house/demo/device">/apiui/latest/house/{houseId}/device</a> — <code>GET /api/latest/house/{houseId}/device</code></li>
  <li><a href="/apiui/latest/ping">/apiui/latest/ping</a> — <code>GET /api/latest/ping</code></li>
</ul>`),
		Script: template.JS(`
window.HataApiUI.setResponse('{"hint":"Pick a page above."}');
window.HataApiUI.setStatus('ready', 'ok');
`),
	})
}

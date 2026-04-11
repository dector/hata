package apiui

import (
	"html/template"
	"net/http"
)

// PingUIHandler serves browser UI page for ping endpoint.
type PingUIHandler struct{}

// NewPingUIHandler creates a new PingUIHandler.
func NewPingUIHandler() *PingUIHandler {
	return &PingUIHandler{}
}

// PingPage serves a UI for GET /api/latest/ping.
func (h *PingUIHandler) PingPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	writePage(w, pageViewData{
		Title:       "Hata API UI — Ping",
		Heading:     "Hata API UI · Ping",
		Description: template.HTML(`Interactive test page for <code>GET /api/latest/ping</code>.`),
		Path:        "/apiui/latest/ping",
		BackHref:    "/apiui",
		Content: template.HTML(`
<div class="row" style="margin-top:0">
  <button id="sendBtn" type="button">Request ping</button>
</div>`),
		CurlSnippet: `curl http://localhost:4501/api/latest/ping`,
		Script: template.JS(`
window.HataApiUI.setResponse('{"hint":"Click request to call /api/latest/ping"}');

document.getElementById('sendBtn')?.addEventListener('click', async () => {
  window.HataApiUI.setStatus('loading...');
  try {
    const res = await fetch('/api/latest/ping');
    const raw = await res.text();
    window.HataApiUI.setResponse(raw);
    window.HataApiUI.setStatus(res.ok ? 'success · ' + res.status : 'error · ' + res.status, res.ok ? 'ok' : 'err');
  } catch (err) {
    window.HataApiUI.setResponse(String(err));
    window.HataApiUI.setStatus('network error', 'err');
  }
});
`),
	})
}

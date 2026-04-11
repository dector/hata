package apiui

import (
	"html/template"
	"net/http"
)

// DeviceUIHandler serves browser UI pages for device endpoints.
type DeviceUIHandler struct{}

// NewDeviceUIHandler creates a new DeviceUIHandler.
func NewDeviceUIHandler() *DeviceUIHandler {
	return &DeviceUIHandler{}
}

// ListPage serves a UI for GET /api/latest/device.
func (h *DeviceUIHandler) ListPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	writePage(w, pageViewData{
		Title:       "Hata API UI — Device list",
		Heading:     "Hata API UI · Device list",
		Description: template.HTML(`Interactive test page for <code>GET /api/latest/device</code>.`),
		Path:        "/apiui/latest/device",
		BackHref:    "/apiui",
		Content: template.HTML(`
<p style="margin-top:0">Uses token from localStorage.</p>
<div class="row">
  <button id="sendBtn" type="button">Request devices</button>
</div>`),
		CurlSnippet: `curl http://localhost:4501/api/latest/device \
  -H "Authorization: Bearer $TOKEN"`,
		Script: template.JS(`
window.HataApiUI.setResponse('{"hint":"Click request to call /api/latest/device"}');

document.getElementById('sendBtn')?.addEventListener('click', async () => {
  const token = window.HataApiUI.loadToken();
  if (!token) {
    window.HataApiUI.setStatus('token is required', 'err');
    return;
  }

  window.HataApiUI.setStatus('loading...');
  try {
    const res = await fetch('/api/latest/device', {
      method: 'GET',
      headers: { Authorization: 'Bearer ' + token },
    });
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

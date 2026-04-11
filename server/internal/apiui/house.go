package apiui

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// HouseUIHandler serves browser UI pages for house endpoints.
type HouseUIHandler struct{}

// NewHouseUIHandler creates a new HouseUIHandler.
func NewHouseUIHandler() *HouseUIHandler {
	return &HouseUIHandler{}
}

// ListPage serves a UI for GET /api/latest/house.
func (h *HouseUIHandler) ListPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	writePage(w, pageViewData{
		Title:       "Hata API UI — House list",
		Heading:     "Hata API UI · House list",
		Description: template.HTML(`Interactive test page for <code>GET /api/latest/house</code>.`),
		Path:        "/apiui/latest/house",
		BackHref:    "/apiui",
		Content: template.HTML(`
<p style="margin-top:0">Uses token from localStorage (or login page).</p>
<div class="row">
  <button id="sendBtn" type="button">Request houses</button>
</div>`),
		CurlSnippet: `curl http://localhost:4501/api/latest/house \
  -H "Authorization: Bearer $TOKEN"`,
		Script: template.JS(`
window.HataApiUI.setResponse('{"hint":"Click request to call /api/latest/house"}');

document.getElementById('sendBtn')?.addEventListener('click', async () => {
  const token = window.HataApiUI.loadToken();
  if (!token) {
    window.HataApiUI.setStatus('token is required', 'err');
    return;
  }

  window.HataApiUI.setStatus('loading...');
  try {
    const res = await fetch('/api/latest/house', {
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

// HouseDevicePage serves a UI for GET /api/latest/house/{houseId}/device.
func (h *HouseUIHandler) HouseDevicePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	houseID := chi.URLParam(r, "houseId")
	content := template.HTML(fmt.Sprintf(`
<label for="houseIdInput">House ID</label>
<input id="houseIdInput" type="text" placeholder="house id" value="%s" />
<div class="row">
  <button id="sendBtn" type="button">Request house devices</button>
</div>`, template.HTMLEscapeString(houseID)))

	writePage(w, pageViewData{
		Title:       "Hata API UI — House devices",
		Heading:     "Hata API UI · House devices",
		Description: template.HTML(`Interactive test page for <code>GET /api/latest/house/{houseId}/device</code>.`),
		Path:        "/apiui/latest/house/{houseId}/device",
		BackHref:    "/apiui",
		Content:     content,
		CurlSnippet: `curl http://localhost:4501/api/latest/house/$HOUSE_ID/device \
  -H "Authorization: Bearer $TOKEN"`,
		Script: template.JS(`
window.HataApiUI.setResponse('{"hint":"Set house id and click request"}');

const houseIdInput = document.getElementById('houseIdInput');

document.getElementById('sendBtn')?.addEventListener('click', async () => {
  const token = window.HataApiUI.loadToken();
  if (!token) {
    window.HataApiUI.setStatus('token is required', 'err');
    return;
  }

  const houseID = (houseIdInput?.value || '').trim();
  if (!houseID) {
    window.HataApiUI.setStatus('houseId is required', 'err');
    return;
  }

  window.HataApiUI.setStatus('loading...');
  try {
    const endpoint = '/api/latest/house/' + encodeURIComponent(houseID) + '/device';
    const res = await fetch(endpoint, {
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

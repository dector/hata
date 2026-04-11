package apiui

import (
	"html/template"
	"net/http"
)

// AuthUIHandler serves browser UI pages for auth endpoints.
type AuthUIHandler struct{}

// NewAuthUIHandler creates a new AuthUIHandler.
func NewAuthUIHandler() *AuthUIHandler {
	return &AuthUIHandler{}
}

// LoginPage serves a UI for POST /api/latest/auth/login.
func (h *AuthUIHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	writePage(w, pageViewData{
		Title:       "Hata API UI — Login",
		Heading:     "Hata API UI · Login",
		Description: template.HTML(`Interactive test page for <code>POST /api/latest/auth/login</code>.`),
		Path:        "/apiui/latest/auth/login",
		BackHref:    "/apiui",
		Content: template.HTML(`
<form id="loginForm">
  <label for="username">Username (email)</label>
  <input id="username" name="username" type="text" autocomplete="username" placeholder="user@example.com" required />

  <label for="password">Password</label>
  <input id="password" name="password" type="password" autocomplete="current-password" placeholder="••••••••" required />

  <div class="row">
    <button type="submit">Send login request</button>
  </div>
</form>`),
		CurlSnippet: `curl -X POST http://localhost:4501/api/latest/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"user@example.com","password":"secret"}'`,
		Script: template.JS(`
const form = document.getElementById('loginForm');
window.HataApiUI.setResponse('{"hint":"Submit credentials to login"}');

form.addEventListener('submit', async (e) => {
  e.preventDefault();
  window.HataApiUI.setStatus('loading...');

  const payload = {
    username: form.username.value,
    password: form.password.value,
  };

  try {
    const res = await fetch('/api/latest/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    });

    const raw = await res.text();
    window.HataApiUI.setResponse(raw);

    if (!res.ok) {
      window.HataApiUI.setStatus('error · ' + res.status, 'err');
      return;
    }

    let parsed;
    try { parsed = JSON.parse(raw); } catch { parsed = null; }
    const token = parsed && parsed.session && parsed.session.token ? parsed.session.token : '';
    if (token) {
      window.HataApiUI.saveToken(token);
      window.HataApiUI.renderTokenStatus();
      window.HataApiUI.setStatus('success · token saved', 'ok');
      return;
    }

    window.HataApiUI.setStatus('success · no token in response', 'ok');
  } catch (err) {
    window.HataApiUI.setResponse(String(err));
    window.HataApiUI.setStatus('network error', 'err');
  }
});
`),
	})
}

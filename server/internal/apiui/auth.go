package apiui

import (
	"net/http"
)

// AuthUIHandler serves minimal browser UI pages for auth endpoints.
type AuthUIHandler struct{}

// NewAuthUIHandler creates a new AuthUIHandler.
func NewAuthUIHandler() *AuthUIHandler {
	return &AuthUIHandler{}
}

// LoginPage serves a small UI for POST /api/latest/auth/login.
func (h *AuthUIHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(loginPageHTML))
}

const loginPageHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>Hata API UI — Login</title>
  <style>
    :root {
      --bg: #0b0f14;
      --panel: #111821;
      --muted: #8ea0b3;
      --text: #e6edf3;
      --accent: #4f8cff;
      --danger: #ff6b6b;
      --ok: #38d39f;
      --border: #1f2a37;
    }
    * { box-sizing: border-box; }
    body {
      margin: 0;
      font-family: Inter, ui-sans-serif, system-ui, -apple-system, Segoe UI, Roboto, Helvetica, Arial, sans-serif;
      background: var(--bg);
      color: var(--text);
      line-height: 1.45;
    }
    .wrap {
      max-width: 860px;
      margin: 32px auto;
      padding: 0 16px;
    }
    .panel {
      background: var(--panel);
      border: 1px solid var(--border);
      border-radius: 12px;
      padding: 18px;
      margin-bottom: 14px;
    }
    h1 { font-size: 20px; margin: 0 0 6px 0; }
    p { margin: 6px 0; color: var(--muted); }
    label { display: block; margin: 12px 0 6px; font-size: 13px; color: var(--muted); }
    input {
      width: 100%;
      padding: 10px 12px;
      border-radius: 8px;
      border: 1px solid var(--border);
      background: #0d141d;
      color: var(--text);
      outline: none;
    }
    input:focus { border-color: var(--accent); }
    .row { display: flex; gap: 10px; align-items: center; margin-top: 14px; }
    button {
      border: 0;
      border-radius: 8px;
      padding: 10px 14px;
      color: white;
      cursor: pointer;
      background: var(--accent);
      font-weight: 600;
    }
    button.secondary { background: #2a3645; }
    .status {
      margin-left: auto;
      font-size: 12px;
      padding: 4px 8px;
      border-radius: 999px;
      border: 1px solid var(--border);
      color: var(--muted);
    }
    .status.ok { color: var(--ok); border-color: rgba(56, 211, 159, 0.5); }
    .status.err { color: var(--danger); border-color: rgba(255, 107, 107, 0.5); }
    pre {
      margin: 0;
      padding: 14px;
      border-radius: 10px;
      background: #0a1017;
      border: 1px solid var(--border);
      overflow: auto;
      font-size: 13px;
    }
    code { color: #c7d2df; }
    .token {
      margin-top: 10px;
      font-size: 13px;
      color: var(--muted);
      word-break: break-all;
    }
  </style>
</head>
<body>
  <main class="wrap">
    <section class="panel">
      <h1>Hata API UI · Login</h1>
      <p>Interactive test page for <code>POST /api/latest/auth/login</code>.</p>
      <p>Page path: <code>/apiui/latest/auth/login</code></p>
    </section>

    <section class="panel">
      <form id="loginForm">
        <label for="username">Username (email)</label>
        <input id="username" name="username" type="text" autocomplete="username" placeholder="user@example.com" required />

        <label for="password">Password</label>
        <input id="password" name="password" type="password" autocomplete="current-password" placeholder="••••••••" required />

        <div class="row">
          <button type="submit">Send login request</button>
          <button id="copyToken" type="button" class="secondary">Copy token</button>
          <span id="status" class="status">idle</span>
        </div>
      </form>
    </section>

    <section class="panel">
      <p style="margin-top:0">Response</p>
      <pre id="response"><code>{
  "hint": "Submit credentials to see response"
}</code></pre>
      <div id="token" class="token"></div>
    </section>
  </main>

  <script>
    const form = document.getElementById('loginForm');
    const statusEl = document.getElementById('status');
    const responseEl = document.getElementById('response');
    const tokenEl = document.getElementById('token');
    const copyBtn = document.getElementById('copyToken');

    let lastToken = '';

    function setStatus(text, kind) {
      statusEl.textContent = text;
      statusEl.classList.remove('ok', 'err');
      if (kind) statusEl.classList.add(kind);
    }

    function setResponseText(text) {
      responseEl.textContent = text;
    }

    form.addEventListener('submit', async (e) => {
      e.preventDefault();
      setStatus('loading...');
      tokenEl.textContent = '';
      lastToken = '';

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
        let data;
        try {
          data = JSON.parse(raw);
          setResponseText(JSON.stringify(data, null, 2));
        } catch {
          setResponseText(raw || '<empty response>');
        }

        if (res.ok) {
          setStatus('success · ' + res.status, 'ok');
          const token = data && data.session && data.session.token ? data.session.token : '';
          if (token) {
            lastToken = token;
            tokenEl.textContent = 'Token: ' + token;
          }
        } else {
          setStatus('error · ' + res.status, 'err');
        }
      } catch (err) {
        setStatus('network error', 'err');
        setResponseText(String(err));
      }
    });

    copyBtn.addEventListener('click', async () => {
      if (!lastToken) {
        setStatus('no token yet', 'err');
        return;
      }
      try {
        await navigator.clipboard.writeText(lastToken);
        setStatus('token copied', 'ok');
      } catch {
        setStatus('copy failed', 'err');
      }
    });
  </script>
</body>
</html>
`

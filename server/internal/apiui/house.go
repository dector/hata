package apiui

import "net/http"

// HouseUIHandler serves minimal browser UI pages for house endpoints.
type HouseUIHandler struct{}

// NewHouseUIHandler creates a new HouseUIHandler.
func NewHouseUIHandler() *HouseUIHandler {
	return &HouseUIHandler{}
}

// ListPage serves a small UI for GET /api/latest/house.
func (h *HouseUIHandler) ListPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(houseListPageHTML))
}

const houseListPageHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>Hata API UI — House list</title>
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
      max-width: 900px;
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
    textarea {
      width: 100%;
      min-height: 96px;
      padding: 10px 12px;
      border-radius: 8px;
      border: 1px solid var(--border);
      background: #0d141d;
      color: var(--text);
      outline: none;
      resize: vertical;
      font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
      font-size: 12px;
    }
    textarea:focus { border-color: var(--accent); }
    .row { display: flex; gap: 10px; align-items: center; margin-top: 14px; flex-wrap: wrap; }
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
    a { color: var(--accent); }
    .back-btn {
      display: inline-block;
      border: 1px solid var(--border);
      border-radius: 8px;
      padding: 6px 10px;
      text-decoration: none;
      background: #0d141d;
      color: var(--text);
      font-size: 13px;
    }
    .back-btn:hover { border-color: var(--accent); }
  </style>
</head>
<body>
  <main class="wrap">
    <section class="panel">
      <h1>Hata API UI · House list</h1>
      <p><a class="back-btn" href="/apiui">← Back</a></p>
      <p>Interactive test page for <code>GET /api/latest/house</code>.</p>
      <p>Page path: <code>/apiui/latest/house</code></p>
      <p>If you need a token, use <a href="/apiui/latest/auth/login">/apiui/latest/auth/login</a>.</p>
    </section>

    <section class="panel">
      <label for="token">Bearer token</label>
      <textarea id="token" spellcheck="false" placeholder="Paste token or login on auth page to auto-fill"></textarea>

      <div class="row">
        <button id="sendBtn" type="button">Request houses</button>
        <button id="saveBtn" type="button" class="secondary">Save token</button>
        <button id="clearBtn" type="button" class="secondary">Clear token</button>
        <span id="status" class="status">idle</span>
      </div>
    </section>

    <section class="panel">
      <p style="margin-top:0">Response</p>
      <pre id="response"><code>{
  "hint": "Click 'Request houses' to call /api/latest/house"
}</code></pre>
    </section>
  </main>

  <script>
    const TOKEN_KEY = 'hata_apiui_token';

    const tokenEl = document.getElementById('token');
    const sendBtn = document.getElementById('sendBtn');
    const saveBtn = document.getElementById('saveBtn');
    const clearBtn = document.getElementById('clearBtn');
    const statusEl = document.getElementById('status');
    const responseEl = document.getElementById('response');

    function setStatus(text, kind) {
      statusEl.textContent = text;
      statusEl.classList.remove('ok', 'err');
      if (kind) statusEl.classList.add(kind);
    }

    function setResponseText(text) {
      responseEl.textContent = text;
    }

    function loadSavedToken() {
      try {
        return localStorage.getItem(TOKEN_KEY) || '';
      } catch {
        return '';
      }
    }

    function saveToken(token) {
      try {
        localStorage.setItem(TOKEN_KEY, token);
      } catch {
        // Ignore localStorage errors
      }
    }

    function clearToken() {
      try {
        localStorage.removeItem(TOKEN_KEY);
      } catch {
        // Ignore localStorage errors
      }
    }

    const savedToken = loadSavedToken();
    if (savedToken) {
      tokenEl.value = savedToken;
      setStatus('saved token loaded', 'ok');
    }

    saveBtn.addEventListener('click', () => {
      const token = tokenEl.value.trim();
      if (!token) {
        setStatus('token is empty', 'err');
        return;
      }
      saveToken(token);
      setStatus('token saved', 'ok');
    });

    clearBtn.addEventListener('click', () => {
      tokenEl.value = '';
      clearToken();
      setStatus('token cleared', 'ok');
    });

    sendBtn.addEventListener('click', async () => {
      const token = tokenEl.value.trim();
      if (!token) {
        setStatus('token is required', 'err');
        return;
      }

      setStatus('loading...');
      saveToken(token);

      try {
        const res = await fetch('/api/latest/house', {
          method: 'GET',
          headers: {
            'Authorization': 'Bearer ' + token,
          },
        });

        const raw = await res.text();
        try {
          const data = JSON.parse(raw);
          setResponseText(JSON.stringify(data, null, 2));
        } catch {
          setResponseText(raw || '<empty response>');
        }

        if (res.ok) {
          setStatus('success · ' + res.status, 'ok');
        } else {
          setStatus('error · ' + res.status, 'err');
        }
      } catch (err) {
        setStatus('network error', 'err');
        setResponseText(String(err));
      }
    });
  </script>
</body>
</html>
`

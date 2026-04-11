package apiui

import (
	"bytes"
	"html/template"
	"net/http"
)

const tokenStorageKey = "hata_apiui_token"

type pageViewData struct {
	Title       string
	Heading     string
	Description template.HTML
	Path        string
	BackHref    string
	Content     template.HTML
	CurlTitle   string
	CurlSnippet string
	Script      template.JS
}

var pageTemplate = template.Must(template.New("apiui-page").Parse(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>{{.Title}}</title>
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
    .wrap { max-width: 920px; margin: 32px auto; padding: 0 16px; }
    .panel {
      background: var(--panel);
      border: 1px solid var(--border);
      border-radius: 12px;
      padding: 18px;
      margin-bottom: 14px;
    }
    h1 { font-size: 20px; margin: 0 0 8px 0; }
    p { margin: 6px 0; color: var(--muted); }
    a { color: var(--accent); }
    ul { margin: 12px 0 0 0; padding-left: 18px; }
    li { margin-bottom: 8px; }
    code { color: #c7d2df; }
    label { display: block; margin: 12px 0 6px; font-size: 13px; color: var(--muted); }
    input, textarea {
      width: 100%;
      padding: 10px 12px;
      border-radius: 8px;
      border: 1px solid var(--border);
      background: #0d141d;
      color: var(--text);
      outline: none;
      font: inherit;
    }
    textarea {
      min-height: 96px;
      resize: vertical;
      font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
      font-size: 12px;
    }
    input:focus, textarea:focus { border-color: var(--accent); }
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
    .back-btn:hover { border-color: var(--accent); text-decoration: none; }
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
    .token-status {
      display: inline-block;
      font-size: 13px;
      border: 1px solid var(--border);
      border-radius: 999px;
      padding: 4px 10px;
      background: #0d141d;
      color: var(--muted);
    }
    .token-status.ok { color: var(--ok); border-color: rgba(56, 211, 159, 0.5); }
    pre {
      margin: 0;
      padding: 14px;
      border-radius: 10px;
      background: #0a1017;
      border: 1px solid var(--border);
      overflow: auto;
      font-size: 13px;
      white-space: pre-wrap;
      word-break: break-word;
    }
  </style>
</head>
<body>
  <main class="wrap">
    <section class="panel">
      <h1>{{.Heading}}</h1>
      <p><a class="back-btn" href="{{.BackHref}}">← Back</a></p>
      <p>{{.Description}}</p>
      <p>Page path: <code>{{.Path}}</code></p>
    </section>

    <section class="panel">
      <p style="margin-top:0">Token status</p>
      <div class="row" style="margin-top:8px">
        <span id="tokenStatus" class="token-status">checking...</span>
        <button id="tokenCopyBtn" type="button" class="secondary">Copy token</button>
        <button id="tokenClearBtn" type="button" class="secondary">Clear token</button>
        <button id="logoutBtn" type="button" class="secondary">Logout session</button>
      </div>
    </section>

    <section class="panel">{{.Content}}</section>

    {{if .CurlSnippet}}
    <section class="panel">
      <p style="margin-top:0">{{if .CurlTitle}}{{.CurlTitle}}{{else}}Copy as curl{{end}}</p>
      <pre id="curlSnippet"><code>{{.CurlSnippet}}</code></pre>
      <div class="row">
        <button id="copyCurlBtn" type="button" class="secondary">Copy curl</button>
      </div>
    </section>
    {{end}}

    <section class="panel">
      <p style="margin-top:0">Response</p>
      <pre id="response"><code>{
  "hint": "Run request to see response"
}</code></pre>
      <div class="row">
        <span id="status" class="status">idle</span>
      </div>
    </section>
  </main>

  <script>
    window.HataApiUI = {
      tokenKey: {{printf "%q" .TokenKey}},
      loadToken() {
        try { return localStorage.getItem(this.tokenKey) || ''; } catch { return ''; }
      },
      saveToken(token) {
        try { localStorage.setItem(this.tokenKey, token); } catch {}
      },
      clearToken() {
        try { localStorage.removeItem(this.tokenKey); } catch {}
      },
      setStatus(text, kind) {
        const el = document.getElementById('status');
        if (!el) return;
        el.textContent = text;
        el.classList.remove('ok', 'err');
        if (kind) el.classList.add(kind);
      },
      setResponse(raw) {
        const el = document.getElementById('response');
        if (!el) return;
        if (typeof raw !== 'string') {
          el.textContent = JSON.stringify(raw, null, 2);
          return;
        }
        try {
          const parsed = JSON.parse(raw);
          el.textContent = JSON.stringify(parsed, null, 2);
        } catch {
          el.textContent = raw || '<empty response>';
        }
      },
      renderTokenStatus() {
        const token = this.loadToken();
        const el = document.getElementById('tokenStatus');
        if (!el) return;
        el.classList.remove('ok');
        if (token) {
          el.textContent = 'present (' + token.length + ' chars)';
          el.classList.add('ok');
        } else {
          el.textContent = 'missing';
        }
      },
      authHeader() {
        const token = this.loadToken();
        return token ? { Authorization: 'Bearer ' + token } : {};
      }
    };

    document.getElementById('tokenCopyBtn')?.addEventListener('click', async () => {
      const token = window.HataApiUI.loadToken();
      if (!token) {
        window.HataApiUI.setStatus('token missing', 'err');
        return;
      }
      try {
        await navigator.clipboard.writeText(token);
        window.HataApiUI.setStatus('token copied', 'ok');
      } catch {
        window.HataApiUI.setStatus('copy failed', 'err');
      }
    });

    document.getElementById('tokenClearBtn')?.addEventListener('click', () => {
      window.HataApiUI.clearToken();
      window.HataApiUI.renderTokenStatus();
      window.HataApiUI.setStatus('token cleared', 'ok');
    });

    document.getElementById('logoutBtn')?.addEventListener('click', async () => {
      const token = window.HataApiUI.loadToken();
      if (!token) {
        window.HataApiUI.setStatus('token missing', 'err');
        return;
      }
      try {
        const res = await fetch('/api/latest/auth/logout', {
          method: 'POST',
          headers: { Authorization: 'Bearer ' + token },
        });
        const raw = await res.text();
        window.HataApiUI.setResponse(raw);

        if (res.status === 404 || res.status === 405) {
          window.HataApiUI.setStatus('logout endpoint unavailable', 'err');
          return;
        }
        if (res.ok) {
          window.HataApiUI.clearToken();
          window.HataApiUI.renderTokenStatus();
          window.HataApiUI.setStatus('logged out', 'ok');
          return;
        }
        window.HataApiUI.setStatus('logout failed · ' + res.status, 'err');
      } catch (err) {
        window.HataApiUI.setResponse(String(err));
        window.HataApiUI.setStatus('network error', 'err');
      }
    });

    document.getElementById('copyCurlBtn')?.addEventListener('click', async () => {
      const snippet = document.getElementById('curlSnippet')?.textContent || '';
      if (!snippet.trim()) return;
      try {
        await navigator.clipboard.writeText(snippet.trim());
        window.HataApiUI.setStatus('curl copied', 'ok');
      } catch {
        window.HataApiUI.setStatus('copy failed', 'err');
      }
    });

    window.HataApiUI.renderTokenStatus();
  </script>
  <script>{{.Script}}</script>
</body>
</html>
`))

func writePage(w http.ResponseWriter, data pageViewData) {
	type pageTemplateData struct {
		Title       string
		Heading     string
		Description template.HTML
		Path        string
		BackHref    string
		Content     template.HTML
		CurlTitle   string
		CurlSnippet string
		Script      template.JS
		TokenKey    string
	}

	tplData := pageTemplateData{
		Title:       data.Title,
		Heading:     data.Heading,
		Description: data.Description,
		Path:        data.Path,
		BackHref:    data.BackHref,
		Content:     data.Content,
		CurlTitle:   data.CurlTitle,
		CurlSnippet: data.CurlSnippet,
		Script:      data.Script,
		TokenKey:    tokenStorageKey,
	}

	var buf bytes.Buffer
	if err := pageTemplate.Execute(&buf, tplData); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(buf.Bytes())
}

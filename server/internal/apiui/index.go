package apiui

import "net/http"

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

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(indexPageHTML))
}

const indexPageHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>Hata API UI</title>
  <style>
    :root {
      --bg: #0b0f14;
      --panel: #111821;
      --muted: #8ea0b3;
      --text: #e6edf3;
      --accent: #4f8cff;
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
    .wrap { max-width: 860px; margin: 32px auto; padding: 0 16px; }
    .panel {
      background: var(--panel);
      border: 1px solid var(--border);
      border-radius: 12px;
      padding: 18px;
      margin-bottom: 14px;
    }
    h1 { font-size: 20px; margin: 0 0 8px 0; }
    p { margin: 0; color: var(--muted); }
    ul { margin: 14px 0 0 0; padding: 0; list-style: none; }
    li {
      padding: 10px 12px;
      border: 1px solid var(--border);
      border-radius: 8px;
      margin-bottom: 8px;
      background: #0d141d;
    }
    a { color: var(--accent); text-decoration: none; }
    a:hover { text-decoration: underline; }
    code { color: #c7d2df; }
    .back-btn {
      display: inline-block;
      border: 1px solid var(--border);
      border-radius: 8px;
      padding: 6px 10px;
      text-decoration: none;
      background: #0d141d;
      color: var(--text);
      font-size: 13px;
      margin-bottom: 8px;
    }
    .back-btn:hover { border-color: var(--accent); text-decoration: none; }
  </style>
</head>
<body>
  <main class="wrap">
    <section class="panel">
      <h1>Hata API UI</h1>
      <a class="back-btn" href="/">← Back</a>
      <p>Available browser test pages under <code>/apiui</code>.</p>
      <ul>
        <li>
          <a href="/apiui/latest/auth/login">/apiui/latest/auth/login</a>
          — interactive login tester for <code>POST /api/latest/auth/login</code>
        </li>
        <li>
          <a href="/apiui/latest/house">/apiui/latest/house</a>
          — interactive house list tester for <code>GET /api/latest/house</code>
        </li>
      </ul>
    </section>
  </main>
</body>
</html>
`

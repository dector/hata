package webauth

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"hata/internal/api"
	"hata/internal/webui"
)

const defaultThenPath = "/me"

// Handler serves browser auth pages.
type Handler struct {
	auth *api.AuthHandler
}

func NewHandler(auth *api.AuthHandler) *Handler {
	return &Handler{auth: auth}
}

func (h *Handler) LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	renderLogin(w, r, loginPageView{
		Then:  sanitizeThen(r.URL.Query().Get("then")),
		Error: "",
	}, http.StatusOK)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	username := strings.TrimSpace(r.FormValue("username"))
	password := r.FormValue("password")
	then := sanitizeThen(r.FormValue("then"))

	result, err := h.auth.Authenticate(r.Context(), username, password)
	if err != nil {
		if errors.Is(err, api.ErrInvalidCredentials) {
			renderLogin(w, r, loginPageView{
				Then:  then,
				Error: "Invalid username or password",
			}, http.StatusUnauthorized)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	api.SetAuthTokenCookie(w, result.Session.Token, result.Session.ValidUntil)
	redirectAfterPost(w, r, then)
}

func (h *Handler) LogoutPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := webui.LogoutPage(webui.LogoutPageData{}).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render logout page", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	api.ClearAuthTokenCookie(w)
	redirectAfterPost(w, r, "/auth/login")
}

type loginPageView struct {
	Then  string
	Error string
}

func renderLogin(w http.ResponseWriter, r *http.Request, view loginPageView, status int) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if err := webui.LoginPage(webui.LoginPageData{Then: view.Then, Error: view.Error}).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render login page", http.StatusInternalServerError)
		return
	}
}

func redirectAfterPost(w http.ResponseWriter, r *http.Request, target string) {
	target = sanitizeThen(target)
	if strings.EqualFold(r.Header.Get("HX-Request"), "true") {
		w.Header().Set("HX-Redirect", target)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	http.Redirect(w, r, target, http.StatusSeeOther)
}

func sanitizeThen(then string) string {
	v := strings.TrimSpace(then)
	if v == "" {
		return defaultThenPath
	}
	if !strings.HasPrefix(v, "/") || strings.HasPrefix(v, "//") {
		return defaultThenPath
	}
	u, err := url.Parse(v)
	if err != nil || u == nil || u.IsAbs() || u.Host != "" {
		return defaultThenPath
	}
	return u.RequestURI()
}


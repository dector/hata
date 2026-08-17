package webauth

import (
	"fmt"
	"hata/internal/webui"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type loginPageView struct {
	Then  string
	Error string
}

func displayNameFromAuth(auth AuthContext) string {
	displayName := strings.TrimSpace(auth.DisplayName)
	if displayName == "" {
		displayName = auth.Username
	}
	return displayName
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
	if isDatastarRequest(r) {
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		_, _ = fmt.Fprintf(w, "window.location.assign(%s);", strconv.Quote(target))
		return
	}
	http.Redirect(w, r, target, http.StatusSeeOther)
}

func sanitizeThen(then string) string {
	v := strings.TrimSpace(then)
	if v == "" {
		return defaultThenPath
	}
	if strings.ContainsAny(v, "\\\r\n") {
		return defaultThenPath
	}
	if !strings.HasPrefix(v, "/") || strings.HasPrefix(v, "//") {
		return defaultThenPath
	}

	u, err := url.ParseRequestURI(v)
	if err != nil || u == nil {
		return defaultThenPath
	}
	if u.IsAbs() || u.Host != "" || u.Scheme != "" {
		return defaultThenPath
	}
	if !strings.HasPrefix(u.Path, "/") || strings.HasPrefix(u.Path, "//") {
		return defaultThenPath
	}
	if _, err := url.ParseQuery(u.RawQuery); err != nil {
		return defaultThenPath
	}
	return u.RequestURI()
}

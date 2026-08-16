package webauth

import (
	"errors"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"hata/internal/api"
	"hata/internal/db"
	"hata/internal/webui"
)

const defaultThenPath = "/me"

// Handler serves browser auth pages.
type Handler struct {
	auth  *api.AuthHandler
	repos db.Repositories
}

func NewHandler(auth *api.AuthHandler, repos db.Repositories) *Handler {
	return &Handler{auth: auth, repos: repos}
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

	api.ClearAuthTokenCookie(w)
	http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	api.ClearAuthTokenCookie(w)
	redirectAfterPost(w, r, "/auth/login")
}

func (h *Handler) MePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	auth, ok := AuthFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := webui.MePage(webui.MePageData{
		UserID:          auth.UserID,
		Username:        auth.Username,
		DisplayName:     displayNameFromAuth(auth),
		HasSessionToken: auth.Token != "",
	}).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render me page", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) AppPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if h.repos == nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	auth, ok := AuthFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	memberships, err := h.repos.HouseRole().ListByUser(r.Context(), auth.UserID)
	if err != nil {
		http.Error(w, "failed to load houses", http.StatusInternalServerError)
		return
	}

	devices, err := h.repos.Device().ListByUser(r.Context(), auth.UserID)
	if err != nil {
		http.Error(w, "failed to load devices", http.StatusInternalServerError)
		return
	}

	devicesByHouse := make(map[string][]webui.AppDeviceData, len(memberships))
	for _, d := range devices {
		devicesByHouse[d.HouseID] = append(devicesByHouse[d.HouseID], webui.AppDeviceData{
			ID:            d.ID,
			Name:          d.Name,
			IntegrationID: d.IntegrationID,
			State:         d.State,
		})
	}

	houses := make([]webui.AppHouseData, 0, len(memberships))
	for _, m := range memberships {
		houses = append(houses, webui.AppHouseData{
			ID:          m.HouseID,
			DisplayName: m.DisplayName,
			Role:        m.Role,
			Devices:     devicesByHouse[m.HouseID],
		})
	}

	sort.Slice(houses, func(i, j int) bool {
		if houses[i].DisplayName == houses[j].DisplayName {
			return houses[i].ID < houses[j].ID
		}
		return houses[i].DisplayName < houses[j].DisplayName
	})

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := webui.AppPage(webui.AppPageData{DisplayName: displayNameFromAuth(auth), Houses: houses}).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render app page", http.StatusInternalServerError)
		return
	}
}

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


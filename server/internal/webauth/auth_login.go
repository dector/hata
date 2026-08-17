package webauth

import (
	"errors"
	"hata/internal/api"
	"net/http"
	"strings"
)

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
			status := http.StatusUnauthorized
			if isDatastarRequest(r) {
				status = http.StatusOK
			}
			renderLogin(w, r, loginPageView{
				Then:  then,
				Error: "Invalid username or password",
			}, status)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	api.SetAuthTokenCookie(w, result.Session.Token, result.Session.ValidUntil)
	redirectAfterPost(w, r, then)
}

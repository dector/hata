package webauth

import (
	"hata/internal/api"
	"net/http"
)

func (h *Handler) LogoutPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	api.ClearAuthTokenCookie(w)
	clearActiveHouseCookie(w)
	http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
}

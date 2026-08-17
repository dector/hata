package webauth

import (
	"hata/internal/api"
	"net/http"
)

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	api.ClearAuthTokenCookie(w)
	redirectAfterPost(w, r, "/auth/login")
}

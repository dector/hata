package webauth

import (
	"hata/internal/webui"
	"net/http"
)

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

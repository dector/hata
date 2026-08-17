package webauth

import (
	"net/http"
)

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

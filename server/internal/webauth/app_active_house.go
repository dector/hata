package webauth

import (
	"net/http"
	"strings"
)

func (h *Handler) SetActiveHouse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	houseID := strings.TrimSpace(r.FormValue("house_id"))
	if houseID == "" {
		http.Error(w, "house id is required", http.StatusBadRequest)
		return
	}

	memberships, err := h.repos.HouseRole().ListByUser(r.Context(), auth.UserID)
	if err != nil {
		http.Error(w, "failed to load houses", http.StatusInternalServerError)
		return
	}
	if findMembership(memberships, houseID) == nil {
		http.NotFound(w, r)
		return
	}

	setActiveHouseCookie(w, houseID)
	then := r.FormValue("then")
	if strings.TrimSpace(then) == "" {
		then = "/app"
	}
	redirectAfterPost(w, r, then)
}

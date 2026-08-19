package webauth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

// HandleHouseExtensionAction handles POST /h/{houseId}/extension/{extensionId}/actions/{action}.
func (h *Handler) HandleHouseExtensionAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	membership := h.authorizedHouseMembership(w, r)
	if membership == nil {
		return
	}
	extensionID := strings.TrimSpace(chi.URLParam(r, "extensionId"))
	action := strings.TrimSpace(chi.URLParam(r, "action"))
	if extensionID == "" || action == "" {
		http.Error(w, "extension ID and action are required", http.StatusBadRequest)
		return
	}

	handler, ok := h.extensions.HouseWebAction(extensionID)
	if !ok {
		http.Error(w, "extension action not found", http.StatusNotFound)
		return
	}
	result, err := handler.HandleHouseWebAction(r.Context(), membership.HouseID, action, r)
	if err != nil {
		http.Error(w, "failed to handle extension action", http.StatusBadGateway)
		return
	}
	status := result.Status
	if status == 0 {
		status = http.StatusOK
	}
	contentType := result.ContentType
	if contentType == "" {
		contentType = "text/html; charset=utf-8"
	}
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(status)
	if result.Render != nil {
		if err := result.Render(r.Context(), w); err != nil {
			fmt.Printf("Error rendering extension action %q/%q for house %q: %v\n", extensionID, action, membership.HouseID, err)
		}
	}
}

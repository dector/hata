package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

// HandleHouseExtensionAction handles POST /api/latest/house/{houseId}/extension/{extensionId}/actions/{action}.
func (h *HouseHandler) HandleHouseExtensionAction(w http.ResponseWriter, r *http.Request) {
	userID, err := userIDFromBearerToken(r, h.repos)
	if err != nil {
		if errors.Is(err, errUnauthorized) {
			WriteError(w, http.StatusUnauthorized, "Unauthorized", "unauthorized")
			return
		}
		fmt.Printf("Error validating token: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return
	}

	houseID := strings.TrimSpace(chi.URLParam(r, "houseId"))
	extensionID := strings.TrimSpace(chi.URLParam(r, "extensionId"))
	action := strings.TrimSpace(chi.URLParam(r, "action"))
	if houseID == "" || extensionID == "" || action == "" {
		WriteError(w, http.StatusBadRequest, "House ID, extension ID, and action are required", "invalid-request")
		return
	}

	allowed, err := h.userCanAccessHouse(r, userID, houseID)
	if err != nil {
		fmt.Printf("Error listing house roles: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return
	}
	if !allowed {
		WriteError(w, http.StatusForbidden, "Forbidden", "forbidden")
		return
	}

	handler, ok := h.houseActionRoutes.HouseAction(extensionID)
	if !ok {
		WriteError(w, http.StatusNotFound, "Extension action not found", "extension-action-not-found")
		return
	}

	result, err := handler.HandleHouseAction(r.Context(), houseID, action, r)
	if err != nil {
		fmt.Printf("Error handling extension action %q/%q for house %q: %v\n", extensionID, action, houseID, err)
		WriteError(w, http.StatusBadGateway, "Failed to handle extension action", "extension-action-failed")
		return
	}
	status := result.Status
	if status == 0 {
		status = http.StatusOK
	}
	WriteJSON(w, status, result.Data)
}

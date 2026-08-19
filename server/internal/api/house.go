package api

import (
	"errors"
	"fmt"
	"net/http"

	"hata/internal/db"
	"hata/internal/extension"
)

// HouseHandler provides house endpoints.
type HouseHandler struct {
	repos      db.Repositories
	extensions *extension.Registry
}

// DefaultHouseExtras lists extension IDs enriched into house responses by default.
var DefaultHouseExtras = []string{
	"hata.ext.weather.v1",
}

// NewHouseHandler creates a new HouseHandler.
func NewHouseHandler(repos db.Repositories) *HouseHandler {
	return NewHouseHandlerWithExtensions(repos, nil)
}

// NewHouseHandlerWithExtensions creates a HouseHandler with extension data.
func NewHouseHandlerWithExtensions(repos db.Repositories, extensions *extension.Registry) *HouseHandler {
	if extensions == nil {
		extensions = extension.NewRegistry()
	}
	return &HouseHandler{repos: repos, extensions: extensions}
}

// SetExtensionRegistry sets generic extension action handlers.
func (h *HouseHandler) SetExtensionRegistry(registry *extension.Registry) {
	if registry == nil {
		registry = extension.NewRegistry()
	}
	h.extensions = registry
}

// List handles GET /api/latest/house
func (h *HouseHandler) List(w http.ResponseWriter, r *http.Request) {
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

	ctx := r.Context()
	memberships, err := h.repos.HouseRole().ListByUser(ctx, userID)
	if err != nil {
		fmt.Printf("Error listing houses: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return
	}

	resp := HouseListResponse{
		Houses: make([]HouseInfo, 0, len(memberships)),
	}
	for _, membership := range memberships {
		houseInfo := HouseInfo{
			ID:          membership.HouseID,
			DisplayName: membership.DisplayName,
			Location:    membership.Location,
			Role:        membership.Role,
		}
		for _, provider := range h.extensions.HouseExtras() {
			if !defaultHouseExtraEnabled(provider.ExtensionID()) {
				continue
			}
			extra, err := provider.HouseExtra(ctx, membership.HouseID)
			if err != nil {
				fmt.Printf("Error loading extension %q for house %q: %v\n", provider.ExtensionID(), membership.HouseID, err)
				continue
			}
			if houseInfo.Extras == nil {
				houseInfo.Extras = map[string]any{}
			}
			houseInfo.Extras[provider.ExtensionID()] = extra
		}
		resp.Houses = append(resp.Houses, houseInfo)
	}

	WriteJSON(w, http.StatusOK, resp)
}

func defaultHouseExtraEnabled(extensionID string) bool {
	for _, enabled := range DefaultHouseExtras {
		if enabled == extensionID {
			return true
		}
	}
	return false
}

func (h *HouseHandler) userCanAccessHouse(r *http.Request, userID int, houseID string) (bool, error) {
	memberships, err := h.repos.HouseRole().ListByUser(r.Context(), userID)
	if err != nil {
		return false, err
	}
	for _, membership := range memberships {
		if membership.HouseID == houseID {
			return true, nil
		}
	}
	return false, nil
}

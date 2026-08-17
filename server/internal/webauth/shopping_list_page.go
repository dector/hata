package webauth

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"hata/internal/db"
	"hata/internal/webui"
)

func (h *Handler) ShoppingListPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	data, ok := h.shoppingListPageData(w, r)
	if !ok {
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := webui.ShoppingListPage(data).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render shopping list page", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) shoppingListPageData(w http.ResponseWriter, r *http.Request) (webui.ShoppingListPageData, bool) {
	if h.repos == nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return webui.ShoppingListPageData{}, false
	}
	auth, ok := AuthFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return webui.ShoppingListPageData{}, false
	}
	listID := strings.TrimSpace(chi.URLParam(r, "listId"))
	if listID == "" {
		http.Error(w, "shopping list id is required", http.StatusBadRequest)
		return webui.ShoppingListPageData{}, false
	}

	memberships, err := h.repos.HouseRole().ListByUser(r.Context(), auth.UserID)
	if err != nil {
		http.Error(w, "failed to load houses", http.StatusInternalServerError)
		return webui.ShoppingListPageData{}, false
	}
	sortMemberships(memberships)
	headerHouses := make([]webui.AppHouseData, 0, len(memberships))
	for _, m := range memberships {
		headerHouses = append(headerHouses, webui.AppHouseData{ID: m.HouseID, DisplayName: m.DisplayName, Role: m.Role})
	}

	activeMembership := activeHouseFromRequest(r, memberships)
	list, membership, err := h.findUserShoppingList(r, listID, activeMembership, memberships)
	if err != nil {
		http.Error(w, "failed to load shopping list", http.StatusInternalServerError)
		return webui.ShoppingListPageData{}, false
	}
	if list == nil || membership == nil {
		http.NotFound(w, r)
		return webui.ShoppingListPageData{}, false
	}
	items, err := h.repos.ShoppingItem().ListByHouseAndListUID(r.Context(), membership.HouseID, list.UID, false)
	if err != nil {
		http.Error(w, "failed to load shopping list items", http.StatusInternalServerError)
		return webui.ShoppingListPageData{}, false
	}
	setActiveHouseCookie(w, membership.HouseID)

	viewItems := make([]webui.ShoppingItemData, 0, len(items))
	for _, item := range items {
		viewItems = append(viewItems, webui.ShoppingItemData{
			ID:      item.UID,
			Name:    item.Name,
			Checked: item.CheckedAt != nil,
		})
	}

	return webui.ShoppingListPageData{
		DisplayName:    displayNameFromAuth(auth),
		HeaderHouses:   headerHouses,
		ActiveHouseID:  membership.HouseID,
		House:          webui.AppHouseData{ID: membership.HouseID, DisplayName: membership.DisplayName, Role: membership.Role},
		List:           webui.ShoppingListData{ID: list.UID, Name: list.Name},
		Items:          viewItems,
		PageURL:        webui.ShoppingListPath(list.UID),
		CanManageHouse: canManageHouse(membership.Role),
	}, true
}

func (h *Handler) userShoppingListForRequest(w http.ResponseWriter, r *http.Request) (*db.ShoppingListData, *db.HouseMembershipData, bool) {
	if h.repos == nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return nil, nil, false
	}
	auth, ok := AuthFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return nil, nil, false
	}
	listID := strings.TrimSpace(chi.URLParam(r, "listId"))
	if listID == "" {
		http.Error(w, "shopping list id is required", http.StatusBadRequest)
		return nil, nil, false
	}
	memberships, err := h.repos.HouseRole().ListByUser(r.Context(), auth.UserID)
	if err != nil {
		http.Error(w, "failed to load houses", http.StatusInternalServerError)
		return nil, nil, false
	}
	activeMembership := activeHouseFromRequest(r, memberships)
	list, membership, err := h.findUserShoppingList(r, listID, activeMembership, memberships)
	if err != nil {
		http.Error(w, "failed to load shopping list", http.StatusInternalServerError)
		return nil, nil, false
	}
	if list == nil || membership == nil {
		http.NotFound(w, r)
		return nil, nil, false
	}
	if !canManageHouse(membership.Role) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return nil, nil, false
	}
	setActiveHouseCookie(w, membership.HouseID)
	return list, membership, true
}

func (h *Handler) findUserShoppingList(r *http.Request, listID string, activeMembership *db.HouseMembershipData, memberships []*db.HouseMembershipData) (*db.ShoppingListData, *db.HouseMembershipData, error) {
	if activeMembership != nil {
		list, err := h.repos.ShoppingList().GetByHouseAndUID(r.Context(), activeMembership.HouseID, listID)
		if err != nil {
			return nil, nil, err
		}
		if list != nil {
			return list, activeMembership, nil
		}
	}
	for _, membership := range memberships {
		if activeMembership != nil && membership.HouseID == activeMembership.HouseID {
			continue
		}
		list, err := h.repos.ShoppingList().GetByHouseAndUID(r.Context(), membership.HouseID, listID)
		if err != nil {
			return nil, nil, err
		}
		if list != nil {
			return list, membership, nil
		}
	}
	return nil, nil, nil
}

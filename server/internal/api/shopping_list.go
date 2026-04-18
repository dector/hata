package api

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"hata/internal/db"

	"github.com/go-chi/chi/v5"
	"github.com/sqids/sqids-go"
)

// ShoppingListHandler provides shopping list and item endpoints.
type ShoppingListHandler struct {
	repos db.Repositories
}

// NewShoppingListHandler creates a new ShoppingListHandler.
func NewShoppingListHandler(repos db.Repositories) *ShoppingListHandler {
	return &ShoppingListHandler{repos: repos}
}

// ListByHouse handles GET /api/latest/house/{houseId}/shopping-list
func (h *ShoppingListHandler) ListByHouse(w http.ResponseWriter, r *http.Request) {
	userID, houseID, ok := h.authorizeHouseAccess(w, r)
	if !ok {
		return
	}
	_ = userID

	lists, err := h.repos.ShoppingList().ListByHouse(r.Context(), houseID)
	if err != nil {
		fmt.Printf("Error listing shopping lists for house %q: %v\n", houseID, err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return
	}

	resp := ShoppingListListResponse{ShoppingLists: make([]ShoppingListInfo, 0, len(lists))}
	for _, list := range lists {
		resp.ShoppingLists = append(resp.ShoppingLists, shoppingListInfoFromData(list))
	}

	WriteJSON(w, http.StatusOK, resp)
}

// GetByHouseAndUID handles GET /api/latest/house/{houseId}/shopping-list/{listId}
func (h *ShoppingListHandler) GetByHouseAndUID(w http.ResponseWriter, r *http.Request) {
	_, houseID, ok := h.authorizeHouseAccess(w, r)
	if !ok {
		return
	}

	listUID := strings.TrimSpace(chi.URLParam(r, "listId"))
	if listUID == "" {
		WriteError(w, http.StatusBadRequest, "List ID is required", "invalid-request")
		return
	}

	list, err := h.repos.ShoppingList().GetByHouseAndUID(r.Context(), houseID, listUID)
	if err != nil {
		fmt.Printf("Error loading shopping list %q for house %q: %v\n", listUID, houseID, err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return
	}
	if list == nil {
		WriteError(w, http.StatusNotFound, "Shopping list not found", "not-found")
		return
	}

	WriteJSON(w, http.StatusOK, ShoppingListGetResponse{ShoppingList: shoppingListInfoFromData(list)})
}

// ListItems handles GET /api/latest/house/{houseId}/shopping-list/{listId}/item
func (h *ShoppingListHandler) ListItems(w http.ResponseWriter, r *http.Request) {
	_, houseID, ok := h.authorizeHouseAccess(w, r)
	if !ok {
		return
	}

	listUID := strings.TrimSpace(chi.URLParam(r, "listId"))
	if listUID == "" {
		WriteError(w, http.StatusBadRequest, "List ID is required", "invalid-request")
		return
	}

	includeDeleted := false
	if raw := strings.TrimSpace(r.URL.Query().Get("include_deleted")); raw != "" {
		parsed, err := strconv.ParseBool(raw)
		if err != nil {
			WriteError(w, http.StatusBadRequest, "include_deleted must be true or false", "invalid-request")
			return
		}
		includeDeleted = parsed
	}

	list, err := h.repos.ShoppingList().GetByHouseAndUID(r.Context(), houseID, listUID)
	if err != nil {
		fmt.Printf("Error loading shopping list %q for house %q: %v\n", listUID, houseID, err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return
	}
	if list == nil {
		WriteError(w, http.StatusNotFound, "Shopping list not found", "not-found")
		return
	}

	items, err := h.repos.ShoppingItem().ListByHouseAndListUID(r.Context(), houseID, listUID, includeDeleted)
	if err != nil {
		fmt.Printf("Error listing shopping items for list %q in house %q: %v\n", listUID, houseID, err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return
	}

	resp := ShoppingItemListResponse{Items: make([]ShoppingItemInfo, 0, len(items))}
	for _, item := range items {
		resp.Items = append(resp.Items, shoppingItemInfoFromData(item))
	}

	WriteJSON(w, http.StatusOK, resp)
}

// CreateItem handles POST /api/latest/house/{houseId}/shopping-list/{listId}/item
func (h *ShoppingListHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	_, houseID, ok := h.authorizeHouseAccess(w, r)
	if !ok {
		return
	}

	listUID := strings.TrimSpace(chi.URLParam(r, "listId"))
	if listUID == "" {
		WriteError(w, http.StatusBadRequest, "List ID is required", "invalid-request")
		return
	}

	list, err := h.repos.ShoppingList().GetByHouseAndUID(r.Context(), houseID, listUID)
	if err != nil {
		fmt.Printf("Error loading shopping list %q for house %q: %v\n", listUID, houseID, err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return
	}
	if list == nil {
		WriteError(w, http.StatusNotFound, "Shopping list not found", "not-found")
		return
	}

	var req ShoppingItemCreateRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body", "invalid-request")
		return
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		WriteError(w, http.StatusBadRequest, "Name is required", "invalid-request")
		return
	}

	itemUID, err := generateUID()
	if err != nil {
		fmt.Printf("Error generating shopping item uid: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return
	}

	item, err := h.repos.ShoppingItem().Create(r.Context(), houseID, listUID, itemUID, name)
	if err != nil {
		fmt.Printf("Error creating shopping item in list %q (house %q): %v\n", listUID, houseID, err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return
	}

	WriteJSON(w, http.StatusCreated, ShoppingItemResponse{Item: shoppingItemInfoFromData(item)})
}

// UpdateItem handles PATCH /api/latest/house/{houseId}/shopping-list/{listId}/item/{itemId}
func (h *ShoppingListHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	_, houseID, ok := h.authorizeHouseAccess(w, r)
	if !ok {
		return
	}

	listUID := strings.TrimSpace(chi.URLParam(r, "listId"))
	if listUID == "" {
		WriteError(w, http.StatusBadRequest, "List ID is required", "invalid-request")
		return
	}
	itemUID := strings.TrimSpace(chi.URLParam(r, "itemId"))
	if itemUID == "" {
		WriteError(w, http.StatusBadRequest, "Item ID is required", "invalid-request")
		return
	}

	var req ShoppingItemUpdateRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body", "invalid-request")
		return
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		WriteError(w, http.StatusBadRequest, "Name is required", "invalid-request")
		return
	}

	item, err := h.repos.ShoppingItem().GetByHouseAndListAndUID(r.Context(), houseID, listUID, itemUID)
	if err != nil {
		fmt.Printf("Error loading shopping item %q in list %q (house %q): %v\n", itemUID, listUID, houseID, err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return
	}
	if item == nil || item.DeletedAt != nil {
		WriteError(w, http.StatusNotFound, "Shopping item not found", "not-found")
		return
	}

	if err := h.repos.ShoppingItem().UpdateName(r.Context(), houseID, listUID, itemUID, name); err != nil {
		fmt.Printf("Error updating shopping item name %q in list %q (house %q): %v\n", itemUID, listUID, houseID, err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return
	}

	updated, err := h.repos.ShoppingItem().GetByHouseAndListAndUID(r.Context(), houseID, listUID, itemUID)
	if err != nil {
		fmt.Printf("Error reloading shopping item %q in list %q (house %q): %v\n", itemUID, listUID, houseID, err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return
	}
	if updated == nil {
		WriteError(w, http.StatusNotFound, "Shopping item not found", "not-found")
		return
	}

	WriteJSON(w, http.StatusOK, ShoppingItemResponse{Item: shoppingItemInfoFromData(updated)})
}

// SetItemChecked handles PATCH /api/latest/house/{houseId}/shopping-list/{listId}/item/{itemId}/check
func (h *ShoppingListHandler) SetItemChecked(w http.ResponseWriter, r *http.Request) {
	userID, houseID, ok := h.authorizeHouseAccess(w, r)
	if !ok {
		return
	}

	listUID := strings.TrimSpace(chi.URLParam(r, "listId"))
	if listUID == "" {
		WriteError(w, http.StatusBadRequest, "List ID is required", "invalid-request")
		return
	}
	itemUID := strings.TrimSpace(chi.URLParam(r, "itemId"))
	if itemUID == "" {
		WriteError(w, http.StatusBadRequest, "Item ID is required", "invalid-request")
		return
	}

	var req ShoppingItemCheckRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body", "invalid-request")
		return
	}

	item, err := h.repos.ShoppingItem().GetByHouseAndListAndUID(r.Context(), houseID, listUID, itemUID)
	if err != nil {
		fmt.Printf("Error loading shopping item %q in list %q (house %q): %v\n", itemUID, listUID, houseID, err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return
	}
	if item == nil || item.DeletedAt != nil {
		WriteError(w, http.StatusNotFound, "Shopping item not found", "not-found")
		return
	}

	var checkedAt *time.Time
	var checkedByUserID *int
	if req.Checked {
		now := time.Now()
		checkedAt = &now
		checkedByUserID = &userID
	}

	if err := h.repos.ShoppingItem().SetChecked(r.Context(), houseID, listUID, itemUID, checkedAt, checkedByUserID); err != nil {
		fmt.Printf("Error setting shopping item checked state %q in list %q (house %q): %v\n", itemUID, listUID, houseID, err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return
	}

	updated, err := h.repos.ShoppingItem().GetByHouseAndListAndUID(r.Context(), houseID, listUID, itemUID)
	if err != nil {
		fmt.Printf("Error reloading shopping item %q in list %q (house %q): %v\n", itemUID, listUID, houseID, err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return
	}
	if updated == nil {
		WriteError(w, http.StatusNotFound, "Shopping item not found", "not-found")
		return
	}

	WriteJSON(w, http.StatusOK, ShoppingItemResponse{Item: shoppingItemInfoFromData(updated)})
}

// DeleteItem handles DELETE /api/latest/house/{houseId}/shopping-list/{listId}/item/{itemId}
func (h *ShoppingListHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	_, houseID, ok := h.authorizeHouseAccess(w, r)
	if !ok {
		return
	}

	listUID := strings.TrimSpace(chi.URLParam(r, "listId"))
	if listUID == "" {
		WriteError(w, http.StatusBadRequest, "List ID is required", "invalid-request")
		return
	}
	itemUID := strings.TrimSpace(chi.URLParam(r, "itemId"))
	if itemUID == "" {
		WriteError(w, http.StatusBadRequest, "Item ID is required", "invalid-request")
		return
	}

	item, err := h.repos.ShoppingItem().GetByHouseAndListAndUID(r.Context(), houseID, listUID, itemUID)
	if err != nil {
		fmt.Printf("Error loading shopping item %q in list %q (house %q): %v\n", itemUID, listUID, houseID, err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return
	}
	if item == nil || item.DeletedAt != nil {
		WriteError(w, http.StatusNotFound, "Shopping item not found", "not-found")
		return
	}

	if err := h.repos.ShoppingItem().SoftDelete(r.Context(), houseID, listUID, itemUID, time.Now()); err != nil {
		fmt.Printf("Error soft-deleting shopping item %q in list %q (house %q): %v\n", itemUID, listUID, houseID, err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ShoppingListHandler) authorizeHouseAccess(w http.ResponseWriter, r *http.Request) (int, string, bool) {
	userID, err := userIDFromBearerToken(r, h.repos)
	if err != nil {
		if errors.Is(err, errUnauthorized) {
			WriteError(w, http.StatusUnauthorized, "Unauthorized", "unauthorized")
			return 0, "", false
		}
		fmt.Printf("Error validating token: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return 0, "", false
	}

	houseID := strings.TrimSpace(chi.URLParam(r, "houseId"))
	if houseID == "" {
		WriteError(w, http.StatusBadRequest, "House ID is required", "invalid-request")
		return 0, "", false
	}

	role, allowed, err := h.userCanAccessShoppingListHouse(r, userID, houseID)
	if err != nil {
		fmt.Printf("Error listing house roles: %v\n", err)
		WriteError(w, http.StatusInternalServerError, "Internal server error", "internal-error")
		return 0, "", false
	}
	if !allowed || role == "guest" {
		WriteError(w, http.StatusForbidden, "Forbidden", "forbidden")
		return 0, "", false
	}

	return userID, houseID, true
}

func (h *ShoppingListHandler) userCanAccessShoppingListHouse(r *http.Request, userID int, houseID string) (string, bool, error) {
	memberships, err := h.repos.HouseRole().ListByUser(r.Context(), userID)
	if err != nil {
		return "", false, err
	}

	for _, membership := range memberships {
		if membership.HouseID != houseID {
			continue
		}
		return strings.ToLower(strings.TrimSpace(membership.Role)), true, nil
	}

	return "", false, nil
}

func shoppingListInfoFromData(list *db.ShoppingListData) ShoppingListInfo {
	return ShoppingListInfo{
		UID:       list.UID,
		Name:      list.Name,
		CreatedAt: list.CreatedAt.Format(time.RFC3339),
		UpdatedAt: list.UpdatedAt.Format(time.RFC3339),
	}
}

func shoppingItemInfoFromData(item *db.ShoppingItemData) ShoppingItemInfo {
	return ShoppingItemInfo{
		UID:             item.UID,
		Name:            item.Name,
		Position:        item.Position,
		CheckedAt:       timePtrToRFC3339(item.CheckedAt),
		CheckedByUserID: item.CheckedByUserID,
		DeletedAt:       timePtrToRFC3339(item.DeletedAt),
		CreatedAt:       item.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       item.UpdatedAt.Format(time.RFC3339),
	}
}

func timePtrToRFC3339(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format(time.RFC3339)
	return &s
}

func generateUID() (string, error) {
	s, err := sqids.New()
	if err != nil {
		return "", fmt.Errorf("failed to create sqids instance: %w", err)
	}

	var buf [8]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return "", fmt.Errorf("failed to read random bytes: %w", err)
	}

	n := binary.BigEndian.Uint64(buf[:])
	if n == 0 {
		n = 1
	}

	uid, err := s.Encode([]uint64{n})
	if err != nil {
		return "", fmt.Errorf("failed to encode uid: %w", err)
	}

	return uid, nil
}

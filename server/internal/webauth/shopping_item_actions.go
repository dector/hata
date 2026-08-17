package webauth

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/sqids/sqids-go"
)

func (h *Handler) AddShoppingItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	list, membership, ok := h.userShoppingListForRequest(w, r)
	if !ok {
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	name := strings.TrimSpace(r.FormValue("name"))
	if name != "" {
		itemUID, err := generateShoppingItemUID()
		if err != nil {
			http.Error(w, "failed to create shopping item", http.StatusInternalServerError)
			return
		}
		if _, err := h.repos.ShoppingItem().Create(r.Context(), membership.HouseID, list.UID, itemUID, name); err != nil {
			http.Error(w, "failed to create shopping item", http.StatusInternalServerError)
			return
		}
	}
	redirectAfterPost(w, r, shoppingListRedirectPath(list.UID))
}

func (h *Handler) SetShoppingItemChecked(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	list, membership, ok := h.userShoppingListForRequest(w, r)
	if !ok {
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	itemUID := strings.TrimSpace(chi.URLParam(r, "itemId"))
	if itemUID == "" {
		http.Error(w, "shopping item id is required", http.StatusBadRequest)
		return
	}
	checked := strings.EqualFold(strings.TrimSpace(r.FormValue("checked")), "true")
	auth, _ := AuthFromContext(r.Context())
	var checkedAt *time.Time
	var checkedByUserID *int
	if checked {
		now := time.Now()
		checkedAt = &now
		checkedByUserID = &auth.UserID
	}
	if err := h.repos.ShoppingItem().SetChecked(r.Context(), membership.HouseID, list.UID, itemUID, checkedAt, checkedByUserID); err != nil {
		http.Error(w, "failed to update shopping item", http.StatusInternalServerError)
		return
	}
	redirectAfterPost(w, r, shoppingListRedirectPath(list.UID))
}

func (h *Handler) RenameShoppingItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	list, membership, ok := h.userShoppingListForRequest(w, r)
	if !ok {
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	itemUID := strings.TrimSpace(chi.URLParam(r, "itemId"))
	name := strings.TrimSpace(r.FormValue("name"))
	if itemUID == "" || name == "" {
		redirectAfterPost(w, r, shoppingListRedirectPath(list.UID))
		return
	}
	if err := h.repos.ShoppingItem().UpdateName(r.Context(), membership.HouseID, list.UID, itemUID, name); err != nil {
		http.Error(w, "failed to rename shopping item", http.StatusInternalServerError)
		return
	}
	redirectAfterPost(w, r, shoppingListRedirectPath(list.UID))
}

func (h *Handler) DeleteShoppingItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	list, membership, ok := h.userShoppingListForRequest(w, r)
	if !ok {
		return
	}
	itemUID := strings.TrimSpace(chi.URLParam(r, "itemId"))
	if itemUID == "" {
		http.Error(w, "shopping item id is required", http.StatusBadRequest)
		return
	}
	if err := h.repos.ShoppingItem().SoftDelete(r.Context(), membership.HouseID, list.UID, itemUID, time.Now()); err != nil {
		http.Error(w, "failed to delete shopping item", http.StatusInternalServerError)
		return
	}
	redirectAfterPost(w, r, shoppingListRedirectPath(list.UID))
}

func shoppingListRedirectPath(listUID string) string {
	return "/sl/" + listUID
}

func generateShoppingItemUID() (string, error) {
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

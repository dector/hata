package api

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"time"

	"hata/internal/db"

	"github.com/sqids/sqids-go"
)

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

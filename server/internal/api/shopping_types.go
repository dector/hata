package api

// ShoppingListInfo contains shopping list details.
type ShoppingListInfo struct {
	UID       string `json:"uid"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// ShoppingListListResponse represents house shopping list endpoint response.
type ShoppingListListResponse struct {
	ShoppingLists []ShoppingListInfo `json:"shoppingLists"`
}

// ShoppingListGetResponse represents single shopping list endpoint response.
type ShoppingListGetResponse struct {
	ShoppingList ShoppingListInfo `json:"shoppingList"`
}

// ShoppingItemInfo contains shopping item details.
type ShoppingItemInfo struct {
	UID             string  `json:"uid"`
	Name            string  `json:"name"`
	Position        int     `json:"position"`
	CheckedAt       *string `json:"checkedAt"`
	CheckedByUserID *int    `json:"checkedByUserId"`
	DeletedAt       *string `json:"deletedAt"`
	CreatedAt       string  `json:"createdAt"`
	UpdatedAt       string  `json:"updatedAt"`
}

// ShoppingItemListResponse represents shopping list items endpoint response.
type ShoppingItemListResponse struct {
	Items []ShoppingItemInfo `json:"items"`
}

// ShoppingItemResponse represents a single shopping item response.
type ShoppingItemResponse struct {
	Item ShoppingItemInfo `json:"item"`
}

// ShoppingItemCreateRequest represents create item payload.
type ShoppingItemCreateRequest struct {
	Name string `json:"name"`
}

// ShoppingItemUpdateRequest represents update item payload.
type ShoppingItemUpdateRequest struct {
	Name string `json:"name"`
}

// ShoppingItemCheckRequest represents check/uncheck payload.
type ShoppingItemCheckRequest struct {
	Checked bool `json:"checked"`
}

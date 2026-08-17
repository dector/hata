package api

// LoginRequest represents the login request payload.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents a successful login response.
type LoginResponse struct {
	Session SessionInfo `json:"session"`
	User    UserInfo    `json:"user"`
}

// SessionInfo contains session details.
type SessionInfo struct {
	Token      string `json:"token"`
	ValidUntil string `json:"validUntil"` // ISO 8601 format
}

// UserInfo contains user details.
type UserInfo struct {
	DisplayName string `json:"displayName"`
}

// PingResponse represents the ping endpoint response.
type PingResponse struct {
	ServerName string `json:"serverName"`
	Version    string `json:"version"`
	APIVersion string `json:"apiVersion"`
}

// HouseInfo contains house details for a user.
type HouseInfo struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	Role        string `json:"role"`
}

// HouseListResponse represents the house list endpoint response.
type HouseListResponse struct {
	Houses []HouseInfo `json:"houses"`
}

// DeviceIntegrationInfo contains integration details.
type DeviceIntegrationInfo struct {
	ID   string         `json:"id"`
	Data map[string]any `json:"data"`
}

// DeviceInfo contains device details.
type DeviceInfo struct {
	ID          string                `json:"id"`
	Name        string                `json:"name"`
	Integration  DeviceIntegrationInfo `json:"integration"`
	State        string                `json:"state"`
	Availability string                `json:"availability"`
}

// HouseRef contains house references for device responses.
type HouseRef struct {
	ID string `json:"id"`
}

// DeviceInfoWithHouse contains device details with house reference.
type DeviceInfoWithHouse struct {
	DeviceInfo
	House HouseRef `json:"house"`
}

// DeviceListResponse represents the device list endpoint response.
type DeviceListResponse struct {
	Devices []DeviceInfoWithHouse `json:"devices"`
}

// DeviceListByHouseResponse represents the house device list endpoint response.
type DeviceListByHouseResponse struct {
	Devices []DeviceInfo `json:"devices"`
}

// DeviceSetStateRequest represents a request to set the desired state for a device.
type DeviceSetStateRequest struct {
	State string `json:"state"`
}

// DeviceSetStateResponse represents the result of setting a device state.
type DeviceSetStateResponse struct {
	DeviceID string   `json:"deviceId"`
	House    HouseRef `json:"house"`
	State    string   `json:"state"`
}

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

// ErrorResponse represents an error response.
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail contains error information.
type ErrorDetail struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

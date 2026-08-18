package repo

import (
	"context"
	"time"
)

// Repositories groups all repository interfaces.
type Repositories interface {
	// KV returns the key-value store repository.
	KV() KVRepository

	// User returns the user repository.
	User() UserRepository

	// Session returns the session repository.
	Session() SessionRepository

	// House returns the house repository.
	House() HouseRepository

	// HouseRole returns the house role repository.
	HouseRole() HouseRoleRepository

	// Device returns the device repository.
	Device() DeviceRepository

	// HouseDiscoveryNetwork returns the house discovery network repository.
	HouseDiscoveryNetwork() HouseDiscoveryNetworkRepository

	// ShoppingList returns the shopping list repository.
	ShoppingList() ShoppingListRepository

	// ShoppingItem returns the shopping item repository.
	ShoppingItem() ShoppingItemRepository
}

// KVRepository provides access to the key-value store.
type KVRepository interface {
	// ListAllKeys returns all keys in the store.
	ListAllKeys(ctx context.Context) ([]string, error)

	// GetByKey retrieves a value by its key.
	// Returns an error if the key doesn't exist.
	GetByKey(ctx context.Context, key string) (string, error)

	// Count returns the total number of key-value pairs.
	Count(ctx context.Context) (int, error)
}

// UserRepository provides access to user data.
type UserRepository interface {
	// GetByID retrieves a user by ID.
	// Returns nil, nil if user doesn't exist.
	GetByID(ctx context.Context, id int) (*UserData, error)

	// GetByUsername retrieves a user by username.
	// Returns nil, nil if user doesn't exist.
	GetByUsername(ctx context.Context, username string) (*UserData, error)

	// Create creates a new user with the given credentials.
	// Returns error if username already exists.
	Create(ctx context.Context, username, passwordHash, displayName string) (*UserData, error)

	// List retrieves users with pagination.
	// Returns a list of users and the total count.
	List(ctx context.Context, limit, offset int) ([]*UserData, int, error)
}

// SessionRepository provides access to session data.
type SessionRepository interface {
	// Create creates a new session for the given user.
	Create(ctx context.Context, userID int, token string, validUntil time.Time) (*SessionData, error)

	// GetByToken retrieves a session by its token.
	// Returns nil, nil if session doesn't exist.
	GetByToken(ctx context.Context, token string) (*SessionData, error)

	// Invalidate marks a session as invalid.
	Invalidate(ctx context.Context, token string) error
}

// HouseRepository provides access to house data.
type HouseRepository interface {
	// Create creates a new house with the given ID and display name.
	Create(ctx context.Context, id, displayName string) (*HouseData, error)

	// GetByID retrieves a house by its ID.
	// Returns nil, nil if house doesn't exist.
	GetByID(ctx context.Context, id string) (*HouseData, error)

	// UpdateLocation updates a house location.
	// Pass nil to clear the location.
	UpdateLocation(ctx context.Context, id string, location *string) error
}

// HouseRoleRepository provides access to house role data.
type HouseRoleRepository interface {
	// Assign assigns a role to a user for a house.
	Assign(ctx context.Context, houseID string, userID int, role string) (*HouseRoleData, error)

	// ListByUser lists houses and roles for the given user.
	ListByUser(ctx context.Context, userID int) ([]*HouseMembershipData, error)
}

// DeviceRepository provides access to device data.
type DeviceRepository interface {
	// Create creates a new device for the given house.
	Create(ctx context.Context, houseID, id, name, integrationID string, integrationData *string, state string) (*DeviceData, error)

	// ListByHouse lists devices for a specific house.
	ListByHouse(ctx context.Context, houseID string) ([]*DeviceData, error)

	// ListByUser lists devices for all houses the user belongs to.
	ListByUser(ctx context.Context, userID int) ([]*DeviceData, error)

	// GetByHouseAndID retrieves a device by house and device ID.
	// Returns nil, nil if device doesn't exist.
	GetByHouseAndID(ctx context.Context, houseID, id string) (*DeviceData, error)

	// UpdateState updates the device state.
	UpdateState(ctx context.Context, houseID, id, state string) error

	// UpdateAvailability updates the device availability.
	UpdateAvailability(ctx context.Context, houseID, id, availability string) error

	// UpdateStatus updates the device state and availability.
	UpdateStatus(ctx context.Context, houseID, id, state, availability string) error

	// UpdateLight updates latest known light settings.
	UpdateLight(ctx context.Context, houseID, id string, brightness *int, colorPreset *string) error

	// ListAll lists all devices.
	ListAll(ctx context.Context) ([]*DeviceData, error)

	// UpdateName updates the device display name.
	UpdateName(ctx context.Context, houseID, id, name string) error
}

// HouseDiscoveryNetworkRepository provides access to house discovery network data.
type HouseDiscoveryNetworkRepository interface {
	// Create creates a discovery network for a house.
	Create(ctx context.Context, houseID, cidr, label string) (*HouseDiscoveryNetworkData, error)

	// ListByHouse lists discovery networks for a house.
	ListByHouse(ctx context.Context, houseID string) ([]*HouseDiscoveryNetworkData, error)

	// DeleteByID deletes a discovery network if it belongs to the house.
	DeleteByID(ctx context.Context, houseID string, id int) error
}

// ShoppingListRepository provides access to shopping list data.
type ShoppingListRepository interface {
	// Create creates a shopping list in a house.
	Create(ctx context.Context, houseID, uid, name string) (*ShoppingListData, error)

	// ListByHouse lists shopping lists for a house.
	ListByHouse(ctx context.Context, houseID string) ([]*ShoppingListData, error)

	// GetByHouseAndUID retrieves a shopping list by house and list UID.
	// Returns nil, nil if list doesn't exist.
	GetByHouseAndUID(ctx context.Context, houseID, uid string) (*ShoppingListData, error)
}

// ShoppingItemRepository provides access to shopping item data.
type ShoppingItemRepository interface {
	// Create creates a shopping item in the specified list and appends it to the end.
	Create(ctx context.Context, houseID, listUID, itemUID, name string) (*ShoppingItemData, error)

	// ListByHouseAndListUID lists shopping items for a list.
	ListByHouseAndListUID(ctx context.Context, houseID, listUID string, includeDeleted bool) ([]*ShoppingItemData, error)

	// GetByHouseAndListAndUID retrieves a shopping item by house/list/item UIDs.
	// Returns nil, nil if item doesn't exist.
	GetByHouseAndListAndUID(ctx context.Context, houseID, listUID, itemUID string) (*ShoppingItemData, error)

	// UpdateName updates shopping item name.
	UpdateName(ctx context.Context, houseID, listUID, itemUID, name string) error

	// SetChecked sets checked state fields.
	// To uncheck, pass checkedAt=nil and checkedByUserID=nil.
	SetChecked(ctx context.Context, houseID, listUID, itemUID string, checkedAt *time.Time, checkedByUserID *int) error

	// SoftDelete marks shopping item as deleted.
	SoftDelete(ctx context.Context, houseID, listUID, itemUID string, deletedAt time.Time) error
}

// UserData represents user information.
type UserData struct {
	ID           int
	Username     string
	PasswordHash string
	DisplayName  string
}

// SessionData represents session information.
type SessionData struct {
	ID           int
	UserID       int
	CreatedAt    time.Time
	ValidUntil   time.Time
	Token        string
	InvalidSince *time.Time
}

// HouseData represents house information.
type HouseData struct {
	ID          string
	DisplayName string
	Location    *string
}

// HouseRoleData represents house role information.
type HouseRoleData struct {
	ID      int
	HouseID string
	UserID  int
	Role    string
}

// HouseMembershipData represents a house membership with role.
type HouseMembershipData struct {
	HouseID     string
	DisplayName string
	Location    *string
	Role        string
}

// DeviceData represents device information.
type DeviceData struct {
	ID               string
	Name             string
	IntegrationID    string
	IntegrationData  *string
	State            string
	Availability     string
	HouseID          string
	LightBrightness  *int
	LightColorPreset *string
}

// HouseDiscoveryNetworkData represents a CIDR configured for house device discovery.
type HouseDiscoveryNetworkData struct {
	ID      int
	HouseID string
	CIDR    string
	Label   string
}

// ShoppingListData represents shopping list information.
type ShoppingListData struct {
	ID        int
	HouseID   string
	UID       string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ShoppingItemData represents shopping item information.
type ShoppingItemData struct {
	ID              int
	ListID          int
	UID             string
	Name            string
	Position        int
	CheckedAt       *time.Time
	CheckedByUserID *int
	DeletedAt       *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

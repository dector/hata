package repo

import (
	"hata/internal/orm"
)

// Repos provides access to all repositories.
type Repos struct {
	kv           *KVRepo
	user         *UserRepo
	session      *SessionRepo
	house        *HouseRepo
	houseRole    *HouseRoleRepo
	device       *DeviceRepo
	shoppingList *ShoppingListRepo
	shoppingItem *ShoppingItemRepo
}

// NewRepos creates a new Repos instance.
func NewRepos(client *orm.Client) *Repos {
	return &Repos{
		kv:           &KVRepo{client: client},
		user:         &UserRepo{client: client},
		session:      &SessionRepo{client: client},
		house:        &HouseRepo{client: client},
		houseRole:    &HouseRoleRepo{client: client},
		device:       &DeviceRepo{client: client},
		shoppingList: &ShoppingListRepo{client: client},
		shoppingItem: &ShoppingItemRepo{client: client},
	}
}

// KV returns the key-value store repository.
func (r *Repos) KV() KVRepository {
	return r.kv
}

// User returns the user repository.
func (r *Repos) User() UserRepository {
	return r.user
}

// Session returns the session repository.
func (r *Repos) Session() SessionRepository {
	return r.session
}

// House returns the house repository.
func (r *Repos) House() HouseRepository {
	return r.house
}

// HouseRole returns the house role repository.
func (r *Repos) HouseRole() HouseRoleRepository {
	return r.houseRole
}

// Device returns the device repository.
func (r *Repos) Device() DeviceRepository {
	return r.device
}

// ShoppingList returns the shopping list repository.
func (r *Repos) ShoppingList() ShoppingListRepository {
	return r.shoppingList
}

// ShoppingItem returns the shopping item repository.
func (r *Repos) ShoppingItem() ShoppingItemRepository {
	return r.shoppingItem
}

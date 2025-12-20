package repo

import (
	"context"
)

// Repositories groups all repository interfaces.
type Repositories interface {
	// KV returns the key-value store repository.
	KV() KVRepository
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

package repo

import (
	"context"
	"fmt"

	"hata/internal/orm"
	"hata/internal/orm/syskv"
)

// KVRepo implements the db.KVRepository interface.
type KVRepo struct {
	client *orm.Client
}

// ListAllKeys returns all keys in the store.
func (k *KVRepo) ListAllKeys(ctx context.Context) ([]string, error) {
	kvs, err := k.client.SysKV.Query().All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying sys__kv: %w", err)
	}

	keys := make([]string, len(kvs))
	for i, kv := range kvs {
		keys[i] = kv.Key
	}

	return keys, nil
}

// GetByKey retrieves a value by its key.
func (k *KVRepo) GetByKey(ctx context.Context, key string) (string, error) {
	kv, err := k.client.SysKV.Query().
		Where(syskv.KeyEQ(key)).
		Only(ctx)

	if err != nil {
		return "", fmt.Errorf("failed querying sys__kv by key %q: %w", key, err)
	}

	return kv.Value, nil
}

// Count returns the total number of key-value pairs.
func (k *KVRepo) Count(ctx context.Context) (int, error) {
	count, err := k.client.SysKV.Query().Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed counting sys__kv: %w", err)
	}

	return count, nil
}

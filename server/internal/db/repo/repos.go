package repo

import (
	"hata/internal/orm"
)

// Repos provides access to all repositories.
type Repos struct {
	kv *KVRepo
}

// NewRepos creates a new Repos instance.
func NewRepos(client *orm.Client) *Repos {
	return &Repos{
		kv: &KVRepo{client: client},
	}
}

// KV returns the key-value store repository.
func (r *Repos) KV() KVRepository {
	return r.kv
}

package repo

import (
	"hata/internal/orm"
)

// Repos provides access to all repositories.
type Repos struct {
	kv      *KVRepo
	user    *UserRepo
	session *SessionRepo
}

// NewRepos creates a new Repos instance.
func NewRepos(client *orm.Client) *Repos {
	return &Repos{
		kv:      &KVRepo{client: client},
		user:    &UserRepo{client: client},
		session: &SessionRepo{client: client},
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

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

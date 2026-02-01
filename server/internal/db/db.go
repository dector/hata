package db

import (
	"context"

	"hata/internal/db/repo"
)

// DB represents the database interface for managing connections and migrations.
type DB interface {
	// Open opens or creates the database.
	// If profilePath is empty, uses 'data/hata.db' as default.
	Open(ctx context.Context, profilePath string) error

	// RunMigrations runs all schema migrations.
	RunMigrations(ctx context.Context) error

	// Repos returns the repositories interface for accessing data.
	Repos() repo.Repositories

	// Close closes the database connection.
	Close() error
}

// Re-export repository interfaces for convenience.
type (
	Repositories      = repo.Repositories
	KVRepository      = repo.KVRepository
	UserRepository    = repo.UserRepository
	SessionRepository = repo.SessionRepository
	UserData          = repo.UserData
	SessionData       = repo.SessionData
)

package db

import (
	"context"
	"database/sql"
	"fmt"

	"hata/internal/db/repo"
	"hata/internal/orm"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"

	_ "github.com/mattn/go-sqlite3"
)

// OpenTestDB opens an in-memory SQLite database for testing.
// Each call creates a new isolated database instance.
// The caller is responsible for calling Close() on the returned DB.
func OpenTestDB(ctx context.Context) (DB, error) {
	// Use :memory: with cache=shared for isolation
	// Each connection gets its own database instance
	dsn := "file::memory:?cache=shared&_fk=1"

	// Open the SQL driver directly so we can apply pragmas
	drv, err := entsql.Open(dialect.SQLite, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed opening in-memory sqlite: %w", err)
	}

	// Apply test-safe pragmas
	if err := useTestSqlitePragmas(drv.DB()); err != nil {
		drv.Close()
		return nil, fmt.Errorf("failed applying test pragmas: %w", err)
	}

	// Create ent client with the configured driver
	client := orm.NewClient(orm.Driver(drv))

	// Create DB implementation
	dbInst := &dbImpl{
		client: client,
		repos:  repo.NewRepos(client),
	}

	// Run migrations
	if err := dbInst.RunMigrations(ctx); err != nil {
		dbInst.Close()
		return nil, fmt.Errorf("failed running migrations: %w", err)
	}

	return dbInst, nil
}

// useTestSqlitePragmas applies test-safe SQLite pragmas to the database connection.
// These pragmas are optimized for in-memory testing and avoid production-only settings.
func useTestSqlitePragmas(conn *sql.DB) error {
	pragmas := []string{
		"PRAGMA foreign_keys = ON",   // Enforce FK constraints in tests
		"PRAGMA synchronous = FULL",  // Data integrity for tests
		"PRAGMA busy_timeout = 1000", // 1 second timeout to avoid flakiness
		// Note: No WAL mode for in-memory databases
		// Note: No page_size/mmap settings - not needed for in-memory
	}

	for _, pragma := range pragmas {
		if _, err := conn.Exec(pragma); err != nil {
			return fmt.Errorf("failed to set pragma %q: %w", pragma, err)
		}
	}

	return nil
}

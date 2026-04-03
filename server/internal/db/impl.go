package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"hata/internal/db/repo"
	"hata/internal/orm"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"

	_ "github.com/mattn/go-sqlite3"
)

// dbImpl is the concrete implementation of the DB interface.
type dbImpl struct {
	client *orm.Client
	repos  *repo.Repos
}

// New creates a new DB instance.
func New() DB {
	return &dbImpl{}
}

// Open opens or creates the database.
func (d *dbImpl) Open(ctx context.Context, profilePath string) error {
	if profilePath == "" {
		// TODO use top-level constant
		profilePath = "data/hata.db"
	}

	if dir := filepath.Dir(profilePath); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("failed creating database directory %q: %w", dir, err)
		}
	}

	// TODO wtf is cache shared?
	dsn := fmt.Sprintf("file:%s?cache=shared&_fk=1", profilePath)

	// Open the SQL driver directly so we can apply pragmas
	drv, err := entsql.Open(dialect.SQLite, dsn)
	if err != nil {
		return fmt.Errorf("failed opening connection to sqlite: %w", err)
	}

	// Apply optimized SQLite pragmas
	if err := UseSqlitePerformancePragmas(drv.DB()); err != nil {
		drv.Close()
		return fmt.Errorf("failed applying pragmas: %w", err)
	}

	// Create ent client with the configured driver
	client := orm.NewClient(orm.Driver(drv))

	d.client = client
	d.repos = repo.NewRepos(client)

	return nil
}

// RunMigrations runs all schema migrations.
func (d *dbImpl) RunMigrations(ctx context.Context) error {
	if d.client == nil {
		return fmt.Errorf("database not opened")
	}

	if err := d.client.Schema.Create(ctx); err != nil {
		return fmt.Errorf("failed creating schema resources: %w", err)
	}

	return nil
}

// Repos returns the repositories interface.
func (d *dbImpl) Repos() repo.Repositories {
	return d.repos
}

// Close closes the database connection.
func (d *dbImpl) Close() error {
	if d.client != nil {
		return d.client.Close()
	}
	return nil
}

// UseSqlitePerformancePragmas applies optimized SQLite pragmas to the database connection.
func UseSqlitePerformancePragmas(conn *sql.DB) error {
	pragmas := []string{
		"PRAGMA journal_mode = WAL",            // Concurrent reads/writes
		"PRAGMA synchronous = NORMAL",          // Faster writes, safe with WAL
		"PRAGMA temp_store = MEMORY",           // Faster temporary storage
		"PRAGMA foreign_keys = ON",             // Enforce referential integrity
		"PRAGMA busy_timeout = 5000",           // Avoid immediate locks (5 seconds)
		"PRAGMA cache_size = -2000",            // ~2MB cache
		"PRAGMA wal_autocheckpoint = 1000",     // Periodic WAL flush
		"PRAGMA mmap_size = 268435456",         // 256MB for faster reads
		"PRAGMA page_size = 4096",              // Modern FS alignment
		"PRAGMA locking_mode = NORMAL",         // Sane file locks
		"PRAGMA secure_delete = OFF",           // Faster deletes
		"PRAGMA journal_size_limit = 67108864", // Cap WAL at 64MB
	}

	for _, pragma := range pragmas {
		if _, err := conn.Exec(pragma); err != nil {
			return fmt.Errorf("failed to set pragma %q: %w", pragma, err)
		}
	}

	return nil
}

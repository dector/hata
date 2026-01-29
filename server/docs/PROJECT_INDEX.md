# Project Index: Hata Server

Generated: 2026-01-30

## 📁 Project Structure

```
hata/server/
├── cmd/hata/          # Main application entry point
├── data/              # SQLite database storage (hata.db)
├── internal/          # Internal packages
│   ├── db/            # Database layer
│   │   ├── repo/      # Repository pattern implementations
│   │   └── schema/    # Ent schema definitions
│   ├── orm/           # Generated ORM code (Ent) [not in repo]
│   └── util/          # Utility packages
└── pkg/               # Public packages
```

## 🚀 Entry Points

- **Main Application**: `cmd/hata/main.go` - HTTP server with chi router, health checks, and database initialization
- **Code Generation**: `internal/generate.go` - Ent ORM code generator (via `go generate`)

## 📦 Core Modules

### Database Layer (`internal/db/`)
- **Path**: `internal/db/db.go`
- **Purpose**: Database interface definition and management
- **Exports**:
  - `DB` interface (Open, RunMigrations, Repos, Close)
  - `Repositories`, `KVRepository` (re-exported from repo package)

### Database Implementation (`internal/db/impl.go`)
- **Path**: `internal/db/impl.go`
- **Purpose**: Concrete implementation of DB interface with SQLite optimizations
- **Key Features**:
  - SQLite connection with WAL mode and performance pragmas
  - Auto-migrations support via Ent
  - Default database path: `data/hata.db`
- **Exports**: `New()` constructor, `UseSqlitePerformancePragmas()` utility

### Repository Layer (`internal/db/repo/`)
- **Interfaces**: `internal/db/repo/interfaces.go`
  - `Repositories` - Repository aggregator
  - `KVRepository` - Key-value store operations (ListAllKeys, GetByKey, Count)
- **Implementation**: `internal/db/repo/repos.go`, `internal/db/repo/kv.go`
  - `Repos` struct aggregates all repositories
  - `KVRepo` implements KVRepository using Ent ORM

### Schema Definitions (`internal/db/schema/`)
- **Path**: `internal/db/schema/syskv.go`
- **Purpose**: Ent schema definition for key-value store
- **Schema**: `SysKV` entity
  - Table: `sys__kv`
  - Fields: `key` (string), `value` (string)

## 🔧 Configuration

- **go.mod**: Go 1.25.5 project with dependencies:
  - `entgo.io/ent` v0.14.5 - ORM framework
  - `github.com/go-chi/chi/v5` v5.2.3 - HTTP router
  - `github.com/mattn/go-sqlite3` v1.14.17 - SQLite driver
  - `ariga.io/atlas` - Database schema migrations
  - Tools: `air` (hot reload), `ent` (code generation)

- **CLAUDE.md**: Project documentation and AI assistant instructions
  - Build command: `ror ai:build`
  - Tech stack overview
  - Development workflow

## 📚 Documentation

- **CLAUDE.md**: Technical overview and development guidelines

## 🧪 Test Coverage

- Unit tests: 0 files
- Integration tests: 0 files
- Total Go files: 8

## 🔗 Key Dependencies

### Direct Dependencies
- **entgo.io/ent** v0.14.5 - Entity framework (ORM) with code generation
- **github.com/go-chi/chi/v5** v5.2.3 - Lightweight HTTP router
- **github.com/mattn/go-sqlite3** v1.14.17 - SQLite database driver

### Development Tools
- **air-verse/air** v1.63.4 - Hot reload for Go applications
- **entgo.io/ent/cmd/ent** - Schema code generator

## 📝 Quick Start

1. **Setup**: Ensure Go 1.25.5+ is installed
2. **Generate ORM code**: `go generate ./...` (generates internal/orm/)
3. **Build**: `ror ai:build` or `go build -o out/hata cmd/hata/main.go`
4. **Run**: `./out/hata` (starts HTTP server on port 8080)
5. **Health check**: `curl http://localhost:8080/health`

## 🏗️ Architecture Patterns

- **Repository Pattern**: Data access abstracted through repository interfaces
- **Code Generation**: ORM code auto-generated from schema definitions using Ent
- **Context-based**: All database operations use `context.Context` for cancellation and timeouts
- **Dependency Injection**: DB instance passed to handlers and services

## 🔍 Key Implementation Details

### SQLite Optimizations
- WAL (Write-Ahead Logging) journal mode for concurrent reads/writes
- Performance pragmas: 2MB cache, 256MB mmap, 4KB page size
- Foreign key enforcement enabled
- 5-second busy timeout to avoid immediate locks

### HTTP Server
- Chi router with middleware: logging, recovery
- Routes: `/` (welcome), `/health` (status check)
- Port: 8080 (or 3000 when running under AIR)
- Graceful database initialization on startup

### Database Workflow
1. Open database (default: `data/hata.db`)
2. Run auto-migrations (creates/updates schema)
3. Access repositories via `DB.Repos()`
4. Use repository methods for type-safe queries

## 📊 Code Statistics

- Total Lines of Code: ~350 (excluding generated ORM)
- Go Files: 8
- Packages: 4 (main, db, repo, schema)
- ORM Generated: Yes (internal/orm/, not tracked in git)

## 🚨 Important Notes

- **Generated Code**: `internal/orm/` is generated and should not be edited manually
- **Migrations**: Auto-run on application startup via `DB.RunMigrations()`
- **Database Location**: Default path is `data/hata.db` (configurable)
- **Code Generation**: Run `go generate ./...` after schema changes

## 📌 Recent Changes (from git)

- `9f5c39e` - ref(server): move server files to server folder
- `55f9fd2` - feat: #vibe WiZ devices status updating
- `8a45e2d` - ref: kotlinize Wiz integration from Cave project (untested)
- `7ff6192` - ref: cut TopBar into components
- `ab4148d` - ref: add Hilt

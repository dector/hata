# Hata - Technical Overview

## Useful commands

- `ror ai:build` - build project.

## Tech Stack

- **Language:** Go 1.25.5
- **Database:** SQLite3 (mattn/go-sqlite3)
- **ORM:** Ent (entgo.io/ent) - code-first schema definition with code generation
- **Migrations:** Ent auto-migrations with Atlas support

## Project Structure

```
hata/
├── cmd/hata/          # Main application entry point
├── internal/          # Internal packages
│   ├── db/            # Database layer
│   │   ├── schema/    # Ent schema definitions
│   │   └── repo/      # Repository pattern implementations
│   └── orm/           # Generated ORM code (Ent)
├── pkg/               # Public packages
├── data/              # SQLite database storage (hata.db)
└── out/               # Build artifacts
```

## Architecture

- **Repository Pattern:** Data access abstracted through repository interfaces
- **Database Layer:** Centralized DB management with connection pooling and migrations
- **Code Generation:** ORM code generated from schema definitions using Ent
- **Context-based:** All database operations use context.Context for cancellation and timeouts

## Database

- SQLite3 database located at `data/hata.db`
- Auto-migrations on startup
- Schema managed through Ent entity definitions in `internal/db/schema/`
- Repository interfaces in `internal/db/repo/` for testability and abstraction

## Development

- Ent code generation: `go generate ./...` (configured in internal/generate.go)
- Schema changes: Modify files in `internal/db/schema/`, then regenerate ORM code

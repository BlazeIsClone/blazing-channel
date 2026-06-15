# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
make build          # compile to bin/api
make run            # go run (no build artifact)
make test           # all tests with -v
make itest          # integration tests for internal/database only (requires Docker)
make watch          # live reload via air (installs air if missing)
make migrate action=up    # apply migrations
make migrate action=down  # rollback migrations
make clean          # remove bin/
```

Run a single test package:

```bash
go test ./post/... -v
go test ./internal/database/... -v   # requires Docker (testcontainers)
```

## Environment

Copy `.env.example` to `.env` and update the variables.

For local dev, start Postgres via Docker Compose (`docker compose up`), then run the API with `make run` or `make watch`.

## Architecture

**Domain packages** (e.g., `post/`): top-level packages, each self-contained with:

- `post.go` — domain types and sentinel errors (`ErrNotFound`)
- `repository.go` — `PostRepo` interface + `PgSQLPostRepo` implementation (raw SQL via `database/sql`)
- `service.go` — `Service` struct; all mutations run inside explicit `BeginTx`/`Commit`/`Rollback` transactions
- `handler.go` — HTTP handler, `Routes()` registers versioned paths (`/api/v1/...`)
- `request.go` — request struct with `validate()` returning `ValidationErrors` (a `map[string]string`)

**Database layer** (`internal/database/`):

- `database.go` — `Service` interface (`Health`, `Close`, `DB() *sql.DB`); singleton connection via pgx/v5 stdlib driver
- `migrate.go` — uses `golang-migrate` with SQL files embedded via `//go:embed migrations/*.sql`

**Shared utilities** (`pkg/api/path.go`): `api.Path(version, path)` builds `/api/<version><path>`.

**Routing convention**: routes registered as `"METHOD /api/v1/resource"` using Go 1.22 method+pattern ServeMux syntax.

**Integration tests** (`internal/database/database_test.go`): spin up a real Postgres container via testcontainers-go — no mocking of the database layer.

## Migrations

SQL files live in `internal/database/migrations/` as `<version>_<name>.up.sql` / `.down.sql` and are embedded at compile time. The standalone `cmd/migrate/main.go` binary accepts `up` or `down` as its first argument.

## Docs

`docs/DOCS.md` is an index of project documentation. `docs/SCHEMA.md` describes the full database schema (tables, indexes, design notes). **Keep both files updated** whenever schema changes or new domain tables are added — propose updates for user approval before writing.

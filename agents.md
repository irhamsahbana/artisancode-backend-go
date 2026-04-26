# AGENTS

Backend guidance is split into the `docs/` folder.

If a task changes backend workflow, coding conventions, or agent behavior, update this `agents.md` file and the relevant `docs/` files in the same task when practical.

## Documentation Index

- [README](./docs/README.md)
- [Architecture](./docs/architecture.md)
- [Data Layers](./docs/data_layers.md)
- [DB & Migration](./docs/db_migration.md)
- [HTTP & Validation](./docs/http_validation.md)
- [Auth Context](./docs/auth_context.md)
- [Localization](./docs/localization.md)
- [Module Integration](./docs/module_integration.md)
- [Error Handling](./docs/error_handling.md)
- [Handler Pattern](./docs/handler_pattern.md)
- [Parameter Convention](./docs/parameter_convention.md)
- [Auth Email Flow](./docs/auth_email_flow.md)


## Backend Reality Check

The active backend production architecture currently centers on these layers:

- HTTP bootstrap lives in `internal/framework/primary/http/`
- HTTP dependency wiring lives in `internal/setup/http_dependency.go`
- Consumer bootstrap lives in `internal/framework/primary/consumer/postgres/`
- Consumer dependency wiring lives in `internal/setup/consumer_dependency.go`
- Active business logic lives in `internal/core/<module>`
- Active HTTP handlers live in `internal/framework/primary/http/<module>`
- Active Postgres repositories live in `internal/framework/secondary/db/postgres/<module>`
- Shared contracts remain in `internal/ports/...`
- `internal/module/` is no longer the dominant production structure; what remains there is template/experimental code

Do not document `internal/module/<module>` as the primary backend pattern when your change touches the active production codebase.

## Shell And Commands

- Prefer the existing `Makefile` as the main command surface when a target exists.
- `Taskfile.yml` is present, but `make` is the clearer documented entry point for routine backend work in this repo.
- Common commands:
  - `make help`
  - `make dev`
  - `go test ./...`
  - `make lint-ci`
  - `make lint-fix`
  - `make ws`
  - `go run ./cmd/bin/main.go consumer`
  - `make scheduler`
  - `make migrate cmd=up`
  - `make create-migration name=create_users_table`
  - `make seed table=rbac`
- There is no dedicated Make test target right now. Use direct `go test` commands for verification, starting with the touched package when possible and escalating to `go test ./...` for broader changes.
- When changing queued email templates or copy, render previews with `go run ./cmd/bin/main.go email-preview`.
- When changing storage upload cleanup or message queue cleanup behavior, use the existing helpers:
  - `make test-storage-upload`
  - `make cleanup-storage-orphans`
  - `make cleanup-message-queue`

## Quick Rules

### Layering

- Handlers only handle HTTP concerns, request validation, mapping, and responses.
- Core contains business rules and orchestration across dependencies.
- Postgres repositories handle SQL queries and scan-result mapping.
- Handlers must not send `restentity` directly into core.
- Core must not import `restentity`.

### File Split Pattern

- Use one file per main operation in handler/core/repository when the module already follows the split-file pattern.
- `handler.go`, `core.go`, and `repo.go` are for struct definitions, config, constructors, and basic registration methods.
- Helpers used by only one operation may stay in that operation file.
- Helpers shared across operations may have their own file when that improves clarity.
- For database migrations, use one migration file per primary schema change unit; related tables should be split into separate sequential timestamped files.

Current active pattern example:

```text
internal/framework/primary/http/attendance/
  handler.go
  get_attendance_logs.go
  get_attendance_log.go
  check_in.go
  check_out.go
  get_attendance_summary_today.go
  get_attendance_policy.go
  get_owner_attendance_dashboard.go
```

### Dependency Injection

- HTTP module wiring is centralized in `internal/setup/http_dependency.go`.
- Consumer wiring is centralized in `internal/setup/consumer_dependency.go`.
- Constructors should use the config-struct pattern.
- Cross-module dependencies must flow through contracts in `internal/ports/...`, not through ad-hoc instantiation inside core or handlers.

### Request Context

- Tenant scope must come from `common.GetUserContext(ctx)`, not from request body or query parameters.
- `restentity` does not store `TenantID`.
- `common.UserContext` currently includes:
  - `UserID`
  - `UserName`
  - `TenantID`
  - `TenantName`
  - `Roles`
  - `CompanyID`
  - `CompanyName`

### Tracing

- Every function that receives `ctx` must start a `tracing.StartSpan`.
- Span format:
  - HTTP handler: `internal:framework:primary:http:<module>:<file>:<function>`
  - core: `internal:core:<module>:<file>:<function>`
  - postgres repository: `internal:framework:secondary:db:postgres:<module>:<file>:<function>`

### Logging

- Use `log.Ctx(ctx)` whenever context is available.
- Parse and validation errors in handlers should use `Warn`.
- Expected business rejections in core should generally use `Warn`.
- Unexpected errors or DB/integration failures should use `Error`.

### Readability Formatting

- Keep backend code vertical when expressions get long. Do not leave wide one-line struct literals, map literals, function calls, or argument lists when they are hard to scan.
- In Postgres repositories, format SQL column lists, `VALUES`, `RETURNING`, `SET`, and multi-condition `WHERE` clauses one item per line when there are multiple fields or predicates.
- Apply this convention consistently across modules, including `internal/core/<module>` and `internal/framework/secondary/db/postgres/<module>`.

### Localization

- Do not parse `Accept-Language` inside feature modules.
- The HTTP app already installs `WithRequestLanguage()` and `LocalizeJSONResponse()`.
- Use `pkg/errmsg` so `message` and `errors` are translated automatically.

### Current Public/Protected Bootstrap Notes

- `internal/framework/primary/http/build.go` registers global middleware, metrics, and calls `setup.HttpDependencies()`.
- User auth public routes are registered under `/users`.
- Most business routes are protected from `app.Group(..., middleware.Auth)` in the setup layer.
- Storage has a mix of protected routes and signed-public-ish routes for private files:
  - protected: `/storage/upload-url`, `/storage/upload`, `/storage/*`
  - signed URL validation: `/storage/private/*`

### Consumers and Jobs

- Active consumer dependencies currently include email subscriptions and export job subscriptions.
- Existing cron tasks include:
  - `cleanup-expired-storage-files`
  - `process-export-jobs`
  - `cleanup-processed-message-queue`

If you add a new consumer or cron task, document it in `docs/module_integration.md` or the relevant domain document.

## PRD Location

- Store backend PRDs in `docs/PRD/`.
- Keep PRD content in Indonesian.
- Treat `docs/PRD/` as local/generated working documentation and do not reference it as mandatory committed engineering documentation.

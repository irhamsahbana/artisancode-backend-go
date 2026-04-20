# Architecture

## Dominant Structure

- Use a per-domain modular pattern at `internal/module/<module_name>`.
- Standard module structure:
  - `repository/` for SQL queries and DB access
  - `core/` for business logic and domain rules
  - `handler/` for HTTP layer

## File Split Pattern (All Layers)

Every module layer (**repository**, **core**, **handler**) must split functions into separate files — one file per main function. The main struct, config, and constructor stay in the base file.

### Repository Layer

```
repository/
  repo.go                # struct, config, constructor only
  get_work_shifts.go     # GetWorkShifts
  get_work_shift.go      # GetWorkShift
  create_work_shift.go   # CreateWorkShift
  update_work_shift.go   # UpdateWorkShift
  delete_work_shift.go   # DeleteWorkShift
```

### Core Layer

```
core/
  core.go                # struct, config, constructor only
  get_work_shifts.go     # GetWorkShifts
  get_work_shift.go      # GetWorkShift
  create_work_shift.go   # CreateWorkShift
  update_work_shift.go   # UpdateWorkShift
  delete_work_shift.go   # DeleteWorkShift
```

### Handler Layer

```
handler/
  handler.go             # struct, config, constructor, Register routes only
  get_work_shifts.go     # getWorkShifts
  get_work_shift.go      # getWorkShift
  create_work_shift.go   # createWorkShift
  update_work_shift.go   # updateWorkShift
  delete_work_shift.go   # deleteWorkShift
```

**Rules:**
- **One main function per file** — no exceptions
- **Helper functions** specific to a main function are placed directly below it in the same file
- The base file (`repo.go`, `core.go`, `handler.go`) contains only: struct definition, config struct, constructor, and (for handler) the `Register` method
- Example: `register.go` contains `RegisterOwner`, `createTenant`, `getOwnerRole`, `createOwnerUser`, while shared email token/template helpers live in dedicated files such as `action_token_helpers.go`

## Ports (Interfaces)

- Interfaces are centralized at `internal/ports/` (not inside modules).
- Structure follows domain/capability grouping:
  - `internal/ports/core/` - service interfaces consumed by primary adapters
  - `internal/ports/primary/` - contracts specific to primary adapters such as HTTP-facing abstractions
  - `internal/ports/secondary/db/` - repository interfaces for driven adapters
  - `internal/ports/secondary/integration/` - integration interfaces for driven adapters
- Modules implement these interfaces, keeping them decoupled and testable.

## Mapper

- Cross-layer mapping is separated under `internal/entity/mapper`.
- Mapper files are split by sub-domain for maintainability.

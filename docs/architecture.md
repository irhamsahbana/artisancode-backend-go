# Architecture

## Dominant Structure

- Use a per-domain modular pattern at `internal/module/<module_name>`.
- Standard module structure:
  - `repository/` for SQL queries and DB access
  - `core/` for business logic and domain rules
  - `handler/` for HTTP layer

## Ports (Interfaces)

- Interfaces are centralized at `internal/ports/` (not inside modules).
- Structure follows domain/capability grouping:
  - `internal/ports/repository/` - repository interfaces
  - `internal/ports/core/` - service interfaces
- Modules implement these interfaces, keeping them decoupled and testable.

## Mapper

- Cross-layer mapping is separated under `internal/entity/mapper`.
- Mapper files are split by sub-domain for maintainability.

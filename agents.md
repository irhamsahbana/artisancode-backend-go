# AGENTS

Technical guidance has been split into the `docs/` folder.

If a task changes backend workflow expectations, coding conventions, or agent behavior, update this `agents.md` file and the relevant `docs/` content in the same task when practical.

## Documentation Index

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

## Quick Rules

### Module Structure
- Each module follows `repository/core/handler` structure
- Ports (interfaces) centralized at `internal/ports/core`, `internal/ports/primary`, `internal/ports/secondary/db`, and `internal/ports/secondary/integration`
- Explicit mappers in `internal/entity/mapper`

### File Split Pattern (All Layers)

Every module layer (**repository**, **core**, **handler**) must split functions into separate files — one file per main function. The main struct, config, and constructor stay in the base file.

#### Repository Layer

```
repository/
  repo.go                # struct, config, constructor only
  get_work_shifts.go     # GetWorkShifts
  get_work_shift.go      # GetWorkShift
  create_work_shift.go   # CreateWorkShift
  update_work_shift.go   # UpdateWorkShift
  delete_work_shift.go   # DeleteWorkShift
```

#### Core Layer

```
core/
  core.go                # struct, config, constructor only
  get_work_shifts.go     # GetWorkShifts
  get_work_shift.go      # GetWorkShift
  create_work_shift.go   # CreateWorkShift
  update_work_shift.go   # UpdateWorkShift
  delete_work_shift.go   # DeleteWorkShift
```

#### Handler Layer

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
- Example: `register.go` contains `RegisterOwner`, `createTenant`, `getOwnerRole`, `createOwnerUser`, `generateAuthTokens`, `sendVerificationEmail`, `buildVerificationEmailBody`

### Core Layer Rules
- Core layer must **never** import `restentity` — it only works with `coreentity`
- Core functions return `coreentity` types; handler maps to `restentity` via mappers
- All core functions with `ctx` parameter must have `tracing.StartSpan`
- Core/helper function span names must follow file path format: `internal:core:<module>:<file_without_extension>:<function>`

### Entity Conventions
- Core entities include `UserCtx common.UserContext` as first field
- Tenant scope uses `tenant_id`
- Foreign keys must validate as `uuidv7`
- `restentity` structs must NOT have `TenantID` - extract from `common.GetUserContext(ctx)`
- Domain string values that already have typed constants in `internal/entity/common/enum.go` must use those constants in code paths, not hardcoded string literals

### Handler to Core Communication
Handler must use mappers to convert `restentity` to `coreentity` before calling core, and map core responses back to `restentity`:

```go
// ✅ Correct - handler uses mapper for input AND output
user := mapper.LoginReqToCore(ctx, *req)
tokens, err := h.core.Login(ctx, user)
resp := mapper.AuthTokensToLoginResp(*tokens)

// ❌ Wrong - handler passes restentity directly to core
resp, err := h.core.Login(ctx, *req)

// ❌ Wrong - core returns restentity
func (c *core) Login(ctx context.Context, user coreentity.User) (*restentity.LoginResp, error)
```

Mapper files are in `internal/entity/mapper/`:
- `user.go` - `LoginReqToCore`, `RegisterReqToCore`, `RefreshTokenReqToCore`, `AuthTokensToLoginResp`, `AuthTokensToRegisterResp`
- `company.go` - `CompanyFromCoreToRest`, `CompanyFromRestCreateToCore`, `CompanyFromRestUpdateToCore`
- `master_employee.go` - `EmployeeFromCoreToRest`, `EmployeeFromRestCreateToCore`, etc.
- `master_org_unit.go` - `OrgUnitFromCoreToRest`, `OrgUnitFromRestCreateToCore`, etc.

### Parameter Convention
- Input parameters use **value structs** (not pointers) — for immutability and nil safety
- Return values use **pointer structs** — to allow nil returns for "not found"
- Exception: use pointer for input only when struct contains `io.Reader`, `*multipart.FileHeader`, etc.
- See [Parameter Convention](./docs/parameter_convention.md) for full details

### Repository
- Use `coreentity` types with filter structs for query parameters
- All methods must use `tracing.StartSpan` for observability
- DB Postgres/helper function span names must follow file path format: `internal:framework:secondary:db:postgres:<module>:<file_without_extension>:<function>`
- Use `r.db.Rebind(query)` consistently for all parameterized queries

### Error Handling
```go
// ✅ Correct
err := someFunction()
if err != nil {
    return nil, err
}

// ❌ Wrong
if err := someFunction(); err != nil {
    return nil, err
}
```

See [Error Handling](./docs/error_handling.md) for status codes, error types, and per-layer patterns.

### Localization
- Request language is already handled centrally. Do not parse `Accept-Language` inside handlers or core logic.
- Fiber middleware `WithRequestLanguage()` stores normalized language in request context.
- Return errors through `pkg/errmsg` so `message` and `errors` can be localized consistently.
- When adding new business error text, register it in `pkg/errmsg/i18n.go` if it must support both `id` and `en`.
- Prefer shared localization flow over hardcoding translated branches in feature modules.

### Logging
- Use `log.Ctx(ctx)` when `ctx` is available (in repository, core, handler layers)
- Use simple, clear messages without `::` prefix
- Use `common.LogKeyPayload` constant for payload key

#### Log Level Rules

| Layer | Situation | Level |
|-------|-----------|-------|
| Handler | Parse/validation error (client fault) | `Warn` |
| Handler | Core/repo error response | `Error` |
| Core | Business rule violation (e.g. duplicate email) | `Warn` |
| Core | Unexpected failure | `Error` |
| Repository | DB query failure | `Error` |

```go
// ✅ Correct - handler parse/validation uses Warn
log.Ctx(ctx).Warn().Err(err).Any(common.LogKeyPayload, req).Msg("Invalid request body")

// ✅ Correct - core business rule uses Warn with specific payload
log.Ctx(ctx).Warn().Any(common.LogKeyPayload, map[string]string{
    "email": req.Email,
}).Msg("Email already registered")

// ✅ Correct - repo DB error uses Error
log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to query employees")

// ❌ Wrong - using Error for client validation failures
log.Ctx(ctx).Error().Err(err).Msg("Failed to parse query params")

// ❌ Wrong - vague or full-struct payload
log.Ctx(ctx).Warn().Any(common.LogKeyPayload, req).Msg("Invalid request")
```

### Tracing
- All functions with `ctx` parameter should have tracing span using `tracing.StartSpan`
- Span names must follow the file path plus function name
- Handler format: `internal:framework:primary:http:<module>:<file_without_extension>:<function>`
- Core format: `internal:core:<module>:<file_without_extension>:<function>`
- DB Postgres format: `internal:framework:secondary:db:postgres:<module>:<file_without_extension>:<function>`

```go
// ✅ Correct - both core and repo have spans with path-based names
func (c *employeeCore) GetEmployees(ctx context.Context, filter coreentity.EmployeeListFilter) ([]coreentity.Employee, int, error) {
    ctx, span := tracing.StartSpan(ctx, "internal:core:employee:get_employees:GetEmployees")
    defer span.End()
    return c.repo.GetEmployees(ctx, filter)
}

func (r *employeeRepo) GetEmployees(ctx context.Context, filter coreentity.EmployeeListFilter) ([]coreentity.Employee, int, error) {
    ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:employee:get_employees:GetEmployees")
    defer span.End()
    // ...
}
```

### Dependency Injection
- Build dependencies in `internal/setup/dependency.go`
- Use config struct pattern for all constructors:

```go
// Config struct pattern
type UserCoreConfig struct {
    Repo       repository.UserRepository
    TokenCache tokencache.TokenCacheContract
}

func NewUserCore(cfg UserCoreConfig) *userCore
```

## User Context & Roles

User supports multiple roles via `user_roles` junction table.

### UserContext Fields
```go
type UserContext struct {
    UserID      string
    UserName    string
    TenantID    string
    TenantName  string
    Roles       []string       // Multiple roles
    CompanyID   *string         // NULL = owner (full access)
    CompanyName *string
}
```

### UserContext Methods
```go
// Check if user is owner (no company restriction)
func (uc UserContext) IsOwner() bool {
    return uc.CompanyID == nil
}

// Check if user has a specific role
func (uc UserContext) HasRole(role string) bool {
    for _, r := range uc.Roles {
        if r == role {
            return true
        }
    }
    return false
}

// Check if user can access a specific company
func (uc UserContext) CanAccessCompany(companyID string) bool {
    if uc.IsOwner() {
        return true  // Owner can access all companies
    }
    return uc.CompanyID != nil && *uc.CompanyID == companyID
}
```

### Access Control Rules
- `company_id = NULL` → Owner (full access to all companies in tenant)
- `company_id = <uuid>` → Company user (restricted to that company only)
- Multiple roles stored in `user_roles` junction table

## Reusable Modules

| Module     | Purpose                              | Usage                          |
| ---------- | ------------------------------------ | ------------------------------ |
| `company`  | Company settings (attendance, leave)   | Per-company configuration      |
| `storage`  | File storage operations              | File upload/delete             |
| `rbac`     | Role-based access control            | Permission checks              |

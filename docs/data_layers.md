# Data Layers

## Entity Layers

- `restentity` for stable HTTP payloads across modules.
- `repoentity` for DB scan results with `db` tags.
- `coreentity` for internal domain models used by core/repository layers.
- Shared domain enums/constants live in `internal/entity/common/enum.go` and should be reused instead of repeating raw string literals such as statuses, sources, and types.

## Core Entity Structure

All `coreentity` structs must include `UserCtx common.UserContext` as the first field:

```go
type Employee struct {
    UserCtx common.UserContext  // Must be first field

    ID            string
    TenantID      string
    EmployeeNo    string
    FullName      string
    // ... other fields
}
```

### Why UserCtx?

The `UserCtx` field serves two purposes:

1. **Ownership Validation**: Resources can be scoped by `tenant_id` OR `user_id` depending on the business logic
2. **Context Propagation**: Audit trail, logging, and tracing can use context info from the request

### UserContext Fields

```go
type UserContext struct {
    UserID      string
    UserName    string
    TenantID    string
    TenantName  string
    Roles       []string   // Multiple roles via user_roles junction table
    CompanyID   *string    // org_unit_id with type='company', NULL = owner (full access)
    CompanyName *string
}
```

### Ownership Patterns

Resources can validate ownership in two ways:

1. **Tenant-based**: `resource.TenantID == userCtx.TenantID`
2. **User-based**: `resource.UserID == userCtx.UserID`

Choose based on the resource's lifecycle:
- Use `tenant_id` for shared resources (org units, job positions)
- Use `user_id` for personal resources (private files, user-specific settings)

## Filter Structs

For repository operations, use dedicated filter structs in `coreentity`:

```go
// Filter for list operations
type EmployeeListFilter struct {
    TenantID  string
    Q         string
    Status    string
    OrgUnitID *string
    Page      int
    Paginate  int
}

// Filter for delete operations
type EmployeeDeleteFilter struct {
    TenantID string
    ID       string
}
```

## Repository Pattern

Repository functions must accept and return `coreentity` types:

```go
// List - accepts filter, returns slice
func (r *repo) GetEmployees(ctx context.Context, filter coreentity.EmployeeListFilter) ([]coreentity.Employee, int, error)

// Get - accepts entity as filter (ID, TenantID), returns entity pointer
func (r *repo) GetEmployee(ctx context.Context, filter coreentity.Employee) (*coreentity.Employee, error)

// Create - accepts entity data, returns entity with ID
func (r *repo) CreateEmployee(ctx context.Context, data coreentity.Employee) (*coreentity.Employee, error)

// Update - accepts entity data
func (r *repo) UpdateEmployee(ctx context.Context, data coreentity.Employee) error

// Delete - accepts delete filter
func (r *repo) DeleteEmployee(ctx context.Context, filter coreentity.EmployeeDeleteFilter) error
```

### Why Filter Structs?

1. **Explicit Intent**: Filter structs clearly indicate which fields are used for filtering
2. **Type Safety**: Compiler catches missing or extra fields
3. **Extensibility**: Easy to add new filter criteria without breaking existing signatures
4. **Consistency**: Same pattern across all repository methods

## Tracing

All repository and core methods should use OpenTelemetry tracing for observability:

```go
import "codebase-app/internal/infrastructure/tracing"

func (r *repo) GetEmployees(ctx context.Context, filter coreentity.EmployeeListFilter) ([]coreentity.Employee, int, error) {
    ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:employee:get_employees:GetEmployees")
    defer span.End()

    // ... method implementation
}
```

Span names must follow the source file path plus function name:

- Handler: `internal:framework:primary:http:<module>:<file_without_extension>:<function>`
- Core: `internal:core:<module>:<file_without_extension>:<function>`
- DB Postgres: `internal:framework:secondary:db:postgres:<module>:<file_without_extension>:<function>`

### Tracing Benefits

1. **Observability**: Track request flow across repository → core → handler
2. **Performance**: Identify slow database queries
3. **Debugging**: Trace errors back to specific operations

## Mapper Pattern

Mappers extract `TenantID` and other user context from `context.Context` automatically:

```go
func EmployeeFromRestCreateToCore(ctx context.Context, req restentity.CreateEmployeeReq) coreentity.Employee {
    uc := common.GetUserContext(ctx)  // Extract from context
    return coreentity.Employee{
        UserCtx:  uc,
        TenantID: uc.TenantID,  // TenantID from context, NOT from request
        // ... other fields
    }
}
```

### Why Extract from Context?

1. **Security**: TenantID comes from authenticated user claims, not from request body
2. **Clean API**: Clients don't need to send tenant_id
3. **Consistency**: All tenant-scoped data uses same source

### restentity Rules

- `restentity` structs must NOT have `TenantID` field
- `restentity` is for HTTP request/response only
- All tenant scoping is handled via context
- If a field represents a constrained domain value backed by shared constants, mapper/core/repository code should use the constant from `common` and only convert to `string` at the boundary when the struct field type still uses `string`

## Mapping Rules

- Use explicit mapper functions in `internal/entity/mapper`.
- Mapper functions accept `context.Context` as the first parameter to enable UserCtx population.
- Avoid implicit mapping to keep transformations clear.

### Mapper Function Signature

```go
func EntityFromRepoToCore(ctx context.Context, item repoentity.Entity) coreentity.Entity
func EntityFromRestCreateToCore(ctx context.Context, req restentity.CreateEntityReq) coreentity.Entity
```

The `ctx` parameter allows mappers to extract user context information when needed for ownership validation.

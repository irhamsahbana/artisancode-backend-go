# Parameter Convention

## Value Struct vs Pointer Struct

### Rule: Input parameters use VALUE struct, return values use POINTER

```go
// ✅ Input = value struct, return = pointer
func (c *core) CreateEmployee(ctx context.Context, data coreentity.Employee) (*coreentity.Employee, error)
func (r *repo) GetEmployees(ctx context.Context, filter coreentity.EmployeeListFilter) ([]coreentity.Employee, int, error)

// ❌ Wrong - input as pointer (unless exception applies)
func (c *core) CreateEmployee(ctx context.Context, data *coreentity.Employee) (*coreentity.Employee, error)
```

### Why Value Struct for Input?

1. **Immutability**: Callee cannot mutate caller's data — no side effects
2. **Nil safety**: No need for `if data == nil` checks at function start
3. **Clarity**: Parameter is always valid and ownable by the function
4. **Performance**: Our structs are small (~10 fields, ~100-200 bytes) — copy cost is negligible

### When to Use Pointer for Input

Use pointer **only** when:

- Struct contains inherently pointer-typed fields (`*multipart.FileHeader`, `io.Reader`)
- Struct is very large (20+ fields or contains large slices/maps)
- Performance profiling shows struct copy is a bottleneck

```go
// ✅ Exception: struct contains *multipart.FileHeader
func (c *core) UploadFile(ctx context.Context, req *coreentity.UploadFileReq) (*coreentity.UploadFileResp, error)
```

## Filter Struct vs Main Entity

### When to Use Filter Struct

Use dedicated filter structs for operations with **different field requirements** than the main entity:

```go
// List filter — only needs filter/pagination fields
type EmployeeListFilter struct {
    TenantID  string
    Q         string
    Status    string
    OrgUnitID *string
    Page      int
    Paginate  int
}

// Delete filter — only needs identity fields
type EmployeeDeleteFilter struct {
    TenantID string
    ID       string
}
```

### When to Use Main Entity as Filter

Use the main entity struct when filtering by identity fields that already exist on the entity:

```go
// Get single — uses entity with ID + TenantID populated
func (r *repo) GetEmployee(ctx context.Context, filter coreentity.Employee) (*coreentity.Employee, error)
```

### Summary

| Operation | Parameter Type | Example |
|-----------|---------------|---------|
| List (with pagination/search) | Dedicated filter struct | `EmployeeListFilter` |
| Get single (by ID) | Main entity struct | `coreentity.Employee` |
| Create | Main entity struct | `coreentity.Employee` |
| Update | Main entity struct | `coreentity.Employee` |
| Delete | Dedicated delete filter struct | `EmployeeDeleteFilter` |

## Mapper Function Signatures

Mapper functions follow consistent naming and parameter patterns:

```go
// Rest → Core (input mapping, accepts ctx for UserContext)
func EntityFromRestCreateToCore(ctx context.Context, req restentity.CreateEntityReq) coreentity.Entity
func EntityFromRestUpdateToCore(ctx context.Context, req restentity.UpdateEntityReq) coreentity.Entity

// Core → Rest (output mapping, no ctx needed)
func EntityFromCoreToRest(item coreentity.Entity) restentity.Entity

// Repo → Core (DB result mapping, accepts ctx for UserContext)
func EntityFromRepoToCore(ctx context.Context, item repoentity.Entity) coreentity.Entity

// Core response → Rest response (output mapping, no ctx needed)
func AuthTokensToLoginResp(tokens coreentity.AuthTokens) restentity.LoginResp
```

### Naming Convention

- Pattern: `{Entity}From{Source}To{Target}`
- For create/update: `{Entity}FromRest{Operation}ToCore`
- Source/Target: `Rest`, `Core`, `Repo`

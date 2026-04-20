# Module Integration

## Module Registration

- Register modules in `internal/setup/http_dependency.go`.
- Initialization order:
  - repository
  - core
  - handler
  - route group

## Route Protection

- Internal routes that require protection should use auth middleware in the module handler.
- Apply `middleware.Auth` on route groups that need authentication.

```go
func (h *handler) Register(router fiber.Router) {
    protected := router.Group("/", middleware.Auth)
    protected.Post("/", h.create)
    protected.Delete("/:id", h.delete)

    // Public routes (no auth)
    router.Get("/public/:id", h.getPublic)
}
```

## Cross-Module Reusability

Modules like `company`, `rbac`, and `storage` are designed to be reusable across other modules. This allows:

- **Configuration**: Fetch company-specific settings via `CompanyRepository`
- **Access Control**: Perform permission checks using `RbacCore`
- **File Operations**: Upload/delete files using `StorageCore`

### Dependency Injection for Cross-Module Usage

When a module needs to use another module's functionality:

```go
// In your module's core constructor
type YourCoreConfig struct {
    Repo    repository.YourRepository
    RbacSvc core.RbacCore              // for permission checks
}

func NewYourCore(cfg YourCoreConfig) *yourCore {
    return &yourCore{
        repo:    cfg.Repo,
        rbacSvc: cfg.RbacSvc,
    }
}
```

### Important Notes

- Always use interface contracts (`ports`) when injecting cross-module dependencies
- Never instantiate module dependencies directly inside handlers or core
- Use centralized ports at `internal/ports/core`, `internal/ports/primary`, `internal/ports/secondary/db`, and `internal/ports/secondary/integration`
- Wire all dependencies in `internal/setup/http_dependency.go`

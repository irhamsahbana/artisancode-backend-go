# Auth Context

## Claims Storage

- Store claims in context as a single `common.UserContext` object.
- Use key `common.UserContextKeyClaims`.

## Standard UserContext Fields

```go
type UserContext struct {
    UserID      string
    UserName    string
    TenantID    string
    TenantName  string
    Roles       []string       // Multiple roles
    CompanyID   *string        // org_unit_id with type='company', NULL = owner (full access)
    CompanyName *string        // org_unit name
}
```

**Note:** `CompanyID` in UserContext is actually the `org_unit_id` of the user's company. The `company_id` column in `users` table references `org_units.id` where `org_units.type = 'company'`.

## UserContext Methods

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

// Check if user can access a specific company (org_unit)
func (uc UserContext) CanAccessCompany(companyID string) bool {
    if uc.IsOwner() {
        return true  // Owner can access all companies
    }
    return uc.CompanyID != nil && *uc.CompanyID == companyID
}
```

## Access Control Rules

- `company_id = NULL` → Owner (full access to all companies in tenant)
- `company_id = <uuid>` → Company user (restricted to that company only)
- Multiple roles stored in `user_roles` junction table

## Company Organization Structure

Companies are stored in `org_units` table with `type = 'company'`:

```
org_units
├── id (UUID)
├── tenant_id (UUID)
├── parent_id (UUID, nullable) - for hierarchy
├── name (VARCHAR)
├── code (VARCHAR, unique per tenant)
├── type (VARCHAR) - 'company', 'division', 'department', 'unit'
├── config (JSONB) - company-specific settings
└── ...
```

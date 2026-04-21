# User Email Tenant Scope Audit

Date: 2026-04-22

## Context

Request: evaluate repository queries and related flows that still treat `users.email` as globally unique instead of tenant-scoped.

Current implementation is inconsistent with tenant-scoped lookup requirements because:

- `FindActiveUserByEmail` only filters by `email` and `deleted_at`.
- Multiple core flows call repository methods that only accept `email`, not `tenant`.
- Database schema still enforces global uniqueness with `UNIQUE (email)`.

Because of that, changing only one SQL query is not sufficient. The issue spans schema, repository contracts, and public auth request shapes.

## Confirmed Global-Email Dependencies

### 1. Schema still enforces global uniqueness

- `db/migrations/20250108070222_create_users_table.sql`
  - `CONSTRAINT users_email_unique UNIQUE (email)`

Impact:

- Cross-tenant duplicate emails are currently impossible at the database level.
- Any tenant-scoped query refactor will remain partial until schema is changed to tenant-aware uniqueness, most likely `UNIQUE (tenant_id, email)` with data migration/backfill review.

### 2. Repository query that directly matches the reported issue

- `internal/framework/secondary/db/postgres/user/find_active_user_by_email.go`
  - Current filter: `WHERE u.email = ? AND u.deleted_at IS NULL`

Risk:

- If duplicate emails per tenant are introduced later, this method can return the wrong tenant's user.

### 3. Repository existence checks still assume email is global

- `internal/framework/secondary/db/postgres/user/register_owner.go`
  - `ExistsActiveUserByEmail`
  - Query: `SELECT id FROM users WHERE email = ? AND deleted_at IS NULL`

- `internal/framework/secondary/db/postgres/user/exists_active_user_by_email_exclude_user.go`
  - `ExistsActiveUserByEmailExcludeUser`
  - Query: `WHERE email = ? AND id <> ? AND deleted_at IS NULL`

Risk:

- Create/update flows will reject an email already used by another tenant, even if tenant-scoped uniqueness is the desired business rule.

## Core Flows Coupled To Global Email Lookup

### A. Public auth flows

- `internal/core/user/forgot_password.go`
  - Calls `FindActiveUserByEmail(ctx, user.Email)`

- `internal/core/user/resend_verification_email.go`
  - Calls `FindActiveUserByEmail(ctx, user.Email)`

- `internal/entity/restentity/user.go`
  - `ForgotPasswordReq` only has `email`
  - `ResendVerificationEmailReq` only has `email`

Impact:

- These flows do not currently carry tenant identity, so they cannot be made tenant-safe by changing repository SQL alone.
- They need a product/API decision first:
  - add `tenant_code` to the public request, or
  - derive tenant from another source with strong guarantees.

### B. Invitation flows

- `internal/core/userinvitation/create_invitation.go`
  - Calls `FindActiveUserByEmail(ctx, data.Email)` even though `data.TenantID` is already known.

- `internal/core/userinvitation/accept_invitation.go`
  - Calls `FindActiveUserByEmail(ctx, item.Email)` even though `item.TenantID` is already known.

Impact:

- These are strong candidates for tenant-scoped refactor because tenant identity is already available in the flow.
- However, the refactor should be paired with schema/contract updates so behavior is consistent end-to-end.

### C. User management flows

- `internal/core/user/create_user.go`
  - Calls `ExistsActiveUserByEmail(ctx, data.Email)` while `data.TenantID` is available.

- `internal/core/user/update_user.go`
  - Calls `ExistsActiveUserByEmailExcludeUser(ctx, data.Email, data.ID)` while `data.TenantID` is available.

- `internal/core/employee/update_employee.go`
  - Calls `userRepo.ExistsActiveUserByEmail(ctx, data.Email)` while `data.TenantID` is available.

- `internal/core/user/register.go`
  - Calls `ExistsActiveUserByEmail(ctx, user.Email)` before tenant is created.

Impact:

- `create_user`, `update_user`, and `employee/update_employee` can become tenant-aware once repository contracts support tenant ID.
- `register.go` is special because owner registration happens before tenant/user creation is complete. Whether this remains globally unique or becomes tenant-scoped depends on the product rule for owner registration.

## Already Tenant-Scoped / Safer References

- `internal/framework/secondary/db/postgres/user/login.go`
  - `FindActiveUserByEmailAndTenant`
  - Filters by both `u.email` and `t.code`

- `internal/framework/secondary/db/postgres/user/find_active_user_by_id_and_tenant.go`
  - Filters by `u.id` and `u.tenant_id`

- `internal/framework/secondary/db/postgres/user/get_user.go`
  - Filters by `u.id` and `u.tenant_id`

- `internal/framework/secondary/db/postgres/user/get_users.go`
  - Filters by `u.tenant_id`

- `internal/framework/secondary/db/postgres/user/update_user_email.go`
  - Updates by `id` and `tenant_id`

- `internal/framework/secondary/db/postgres/user/update_user_password.go`
  - Updates by `id` and `tenant_id`

- `internal/framework/secondary/db/postgres/user/delete_user.go`
  - Deletes by `id` and `tenant_id`

## Recommended Refactor Order

1. Decide target business rule for user email:
   - global unique across all tenants, or
   - unique per tenant.

2. If target is tenant-scoped:
   - replace schema uniqueness from `UNIQUE (email)` to tenant-aware uniqueness,
   - add/rename repository methods so tenant is explicit in method signature,
   - update invitation/user-management flows that already know `tenant_id`,
   - redesign public auth payloads (`forgot password`, `resend verification`) to include tenant identity.

3. After contract changes, retire ambiguous methods:
   - `FindActiveUserByEmail`
   - `ExistsActiveUserByEmail`
   - `ExistsActiveUserByEmailExcludeUser`

## Files That Need Evaluation

- `db/migrations/20250108070222_create_users_table.sql`
- `internal/ports/secondary/db/user.go`
- `internal/framework/secondary/db/postgres/user/find_active_user_by_email.go`
- `internal/framework/secondary/db/postgres/user/register_owner.go`
- `internal/framework/secondary/db/postgres/user/exists_active_user_by_email_exclude_user.go`
- `internal/core/user/forgot_password.go`
- `internal/core/user/resend_verification_email.go`
- `internal/core/user/create_user.go`
- `internal/core/user/update_user.go`
- `internal/core/user/register.go`
- `internal/core/employee/update_employee.go`
- `internal/core/userinvitation/create_invitation.go`
- `internal/core/userinvitation/accept_invitation.go`
- `internal/entity/restentity/user.go`
- `internal/entity/mapper/user.go`

## Conclusion

There are confirmed queries and flows that still rely on global email lookup. The originally reported file is one of them, but it is not an isolated bug. A safe fix requires coordinated changes across:

- database uniqueness constraint,
- repository method contracts,
- core use cases,
- and public auth request payloads.

# Categories

## Purpose

- `categories` is the general-purpose master-data table for simple reusable category sets.
- Use it for lightweight lookup data that may need hierarchy, tenant scope, and optional org-unit scope.
- Do not use it as a catch-all table for domain entities with their own workflow or rich relations.

## When To Use

Use `categories` for cases such as:

- expense categories
- leave reasons
- simple document tags
- lightweight dropdown options with parent-child structure

Create a dedicated table instead when the data:

- has its own workflow or approval logic
- has many domain-specific relations
- needs behavior beyond simple lookup and hierarchy

## Schema Intent

Main fields in [`db/migrations/20250108070016_create_categories_table.sql`](./../db/migrations/20250108070016_create_categories_table.sql):

- `tenant_id`: required tenant ownership for every row
- `org_unit_id`: optional narrower scope inside the tenant
- `parent_id`: self-reference for simple category trees
- `group_key`: logical grouping for reusable category sets, for example `expense_category` or `leave_reason`
- `code`: stable internal identifier for seeds, filters, and business rules
- `name`: display label that may change over time
- `notes`: optional human-readable explanation
- `sort_order`: ordering hint for lists
- `is_active`: soft activation flag
- `metadata`: small structured attributes when needed

## Scope Rules

- `tenant_id` must always be present.
- `org_unit_id IS NULL` means the category is available at tenant level.
- `org_unit_id IS NOT NULL` means the category is scoped to one specific org unit inside that tenant.
- `parent_id` should only point to another category in the same logical group.

## Code Usage Rules

- Use `group_key` to separate unrelated category families in the same table.
- Use `code` for internal logic, seeds, filtering, and integration mappings.
- Do not use `name` as an identifier in code because display wording can change.
- Treat `metadata` as lightweight extension data, not as a place to move core domain logic.

## Query And Uniqueness Rules

- Tenant-level categories are unique by `(tenant_id, group_key, code)` when `org_unit_id IS NULL`.
- Org-unit-level categories are unique by `(tenant_id, org_unit_id, group_key, code)` when `org_unit_id IS NOT NULL`.
- Query by `tenant_id` first, then narrow by `group_key`, and then by `org_unit_id` when needed.

## Migration Position

- Create `categories` after `tenants` and `org_units` because it references both.
- Keep future schema changes for `categories` in separate migration files rather than rewriting the initial table migration after rollout.

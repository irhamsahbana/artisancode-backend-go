# DB & Migration

## Tenant Scope

- All multi-tenant scope must use `tenant_id`.
- Avoid using `company_id` in new designs.

## Timestamp Columns

- `created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP`
- `updated_at TIMESTAMP WITH TIME ZONE`
- `deleted_at TIMESTAMP WITH TIME ZONE`

## Naming Migration

- Use file naming pattern `YYYYMMDDHHMMSS_<name>.sql`.
- Keep one migration file for one schema change unit so rollout and rollback stay explicit.
- If a change creates related tables, split them into sequential migration files based on dependency order.

## Local Restore

- For local PostgreSQL rehydration from `./backups/`, use `make restore file=<backup-file>`.
- The repo wrapper runs `pg_restore --clean`, so treat it as a destructive local-data workflow and only use it when the task explicitly requires a restore.

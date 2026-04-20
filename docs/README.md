# Documentation Index

This documentation splits engineering guidance into smaller files for more efficient consumption by agents and developers.

When a backend task changes workflow expectations, coding conventions, or agent operating instructions, update the relevant `agents.md` file alongside the affected documentation when practical.

## Document List

- [Architecture](./architecture.md) - Module structure, ports, mapper pattern
- [Data Layers](./data_layers.md) - Entity layers (rest/core/repo), filter structs, tracing
- [DB & Migration](./db_migration.md) - Tenant scope, timestamps, migration naming
- [HTTP & Validation](./http_validation.md) - Request structs, validation tags, FK rules
- [Auth Context](./auth_context.md) - UserContext, access control, company structure
- [Module Integration](./module_integration.md) - DI, route protection, cross-module usage
- [Error Handling](./error_handling.md) - Error types, status codes, per-layer patterns
- [Handler Pattern](./handler_pattern.md) - CRUD handler boilerplate, pagination, log levels
- [Parameter Convention](./parameter_convention.md) - Value vs pointer, filter structs, mapper naming
- [Storage Upload Flow](./storage_upload_flow.md) - Presigned upload lifecycle, storage file linking, orphan cleanup
- [Mobile Attendance V1](./mobile_attendance_v1.md) - Mobile readiness scope, endpoint contracts, and backlog
- [Auth Email Flow](./auth_email_flow.md) - Verification/reset token flow, rate limiting, templates, and preview command

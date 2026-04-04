# Storage Upload Flow

## Overview

The storage upload flow uses two tables:

- `storage_files` stores file metadata and lifecycle state.
- `storage_file_links` stores the relationship between a stored file and a business resource.

This keeps upload lifecycle concerns separate from domain resource linking.

## Table Responsibilities

### `storage_files`

This table stores:

- uploader and tenant ownership
- `object_key`
- `content_type`
- `size_bytes`
- lifecycle state such as `pending`, `attached`, `deleted`, `failed`
- expiration and cleanup timestamps

### `storage_file_links`

This table stores:

- the stored file being linked via `storage_file_id`
- target resource identity via `resource_type` and `resource_id`
- semantic usage via `field_name`
- ordering via `sort_order`

Example:

- `resource_type = attendance_log`
- `resource_id = <attendance_log_id>`
- `field_name = selfie`
- `sort_order = 1`

## Upload Flow

### 1. Client requests upload URL

Client calls:

- `POST /storage/upload-url`

Backend behavior:

- generates the final `object_key`
- creates a `storage_files` row with `status = pending`
- returns presigned upload URL and `file_id`

### 2. Client uploads binary file to object storage

Client uploads directly to the returned presigned `PUT` URL.

Important:

- use the raw file body
- do not use multipart form-data
- send the same `Content-Type` used during presign

### 3. Client submits business request

Example:

- `POST /attendance-logs/check-in`

Request includes:

- `selfie_file_id`

### 4. Domain record is created

Example:

- backend creates `attendance_logs` row

### 5. File is linked to the domain record

Backend creates a `storage_file_links` row such as:

- `storage_file_id = <file_id>`
- `resource_type = attendance_log`
- `resource_id = <attendance_log_id>`
- `field_name = selfie`
- `sort_order = 1`

### 6. File lifecycle is updated

After link creation succeeds:

- `storage_files.status` is updated from `pending` to `attached`

## Cleanup Orphan Files

Orphan files are files that:

- are still `pending`
- are expired
- do not have any row in `storage_file_links`

Cleanup process:

1. query expired `pending` rows in `storage_files`
2. ensure no matching rows exist in `storage_file_links`
3. delete object from object storage
4. mark `storage_files.status = deleted`
5. set `deleted_at`

Command:

```bash
go run ./cmd/bin/main.go cronjob --task=cleanup-expired-storage-files --limit=100
```

Or with Make:

```bash
make cleanup-storage-orphans
```

## Private File Access

Private files use their full object key in the route path.

Example:

- object key: `private/tenants/<tenant-id>/attendance-face/users/<user-id>/<file>.jpg`
- access route: `GET /storage/private/tenants/<tenant-id>/attendance-face/users/<user-id>/<file>.jpg?...`

The backend resolves the full object key from the wildcard route and validates the file metadata before creating the presigned download URL.

## Delete Behavior

Delete requests also use the full object key in the route path.

Example:

- `DELETE /storage/private/tenants/<tenant-id>/attendance-face/users/<user-id>/<file>.jpg`

Only files that still belong to the authenticated tenant and are still in `pending` status can be deleted. When deletion succeeds, the object is removed from object storage and `storage_files` is marked `deleted`.

## Why This Design

This design supports:

- temporary uploads before business form submission
- orphan cleanup
- one resource with many files
- one field with many files using `sort_order`
- future reuse of the same storage module across attendance, reimbursement, profile, and documents

## Attendance Example

For attendance photo proof:

1. create upload URL
2. upload selfie to object storage
3. submit check-in with `selfie_file_id`
4. create attendance log
5. create `storage_file_links` row with:
   - `resource_type = attendance_log`
   - `field_name = selfie`
6. mark file `attached`

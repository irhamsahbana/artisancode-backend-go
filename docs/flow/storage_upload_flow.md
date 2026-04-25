# Storage Upload Flow

## Overview

Storage currently supports two upload paths:

1. presigned upload via `POST /storage/upload-url`
2. direct multipart upload via `POST /storage/upload`

For attendance selfies and mobile flows, the recommended path remains presigned upload because the backend tracks the file lifecycle before the file is attached to a business resource.

## Main Tables

### `storage_files`

Stores file metadata and upload lifecycle state:

- tenant ownership
- creator/uploader
- object key
- original filename
- content type
- file status such as `pending`, `attached`, and `deleted`

### `storage_file_links`

Stores the relationship between files and business resources:

- `storage_file_id`
- `resource_type`
- `resource_id`
- `field_name`
- `sort_order`

Attendance selfie example:

- `resource_type = attendance_log`
- `field_name = selfie`

## Recommended Presigned Flow

### 1. Client requests an upload URL

Request:

- `POST /storage/upload-url`

Minimum body:

- `filename`
- `content_type`
- `folder`
- `is_public`

The backend will:

- generate the final object key
- create a `storage_files` row with initial status `pending`
- return `file_id`, `object_key`, `upload_url`, `method`, and `headers`

### 2. Client uploads the binary to object storage

The client uploads the file to the presigned URL.

Important:

- use the HTTP method and headers returned by the backend
- the content type must match the one used during presign

### 3. Client submits the business request

Attendance examples:

- `POST /attendance-logs/check-in`
- `POST /attendance-logs/check-out`

Current attendance payloads must include:

- `selfie_file_id`

### 4. Business record is created

Example:

- the backend creates `attendance_logs`

### 5. File is linked to the resource

The backend creates a row in `storage_file_links`.

### 6. File is marked as attached

After the link succeeds, the backend calls `MarkFileAttached(...)` and moves the file status from `pending` to `attached`.

## Direct Multipart Flow

Available route:

- `POST /storage/upload`

This flow pushes the file directly to the storage integration without the `presign -> attach` lifecycle used by attendance.

It is usually suitable for admin or internal use cases that only need a simple upload. If the use case requires safely attaching files to a business resource and cleaning up orphan files, the presigned flow is safer.

## Private File Access

Private files are accessed through:

- `GET /storage/private/*`

This route uses signed URL validation middleware, then the backend maps the wildcard path to the real private object key before generating the file URL.

## Delete Behavior

File deletion uses:

- `DELETE /storage/*`

The delete handler currently passes the object key/path to core. Files may only be deleted when they satisfy the ownership and state rules enforced by core and repository logic.

## Orphan Cleanup

Active cron task:

```bash
go run ./cmd/bin/main.go cronjob --task=cleanup-expired-storage-files --limit=100
```

Or via Make:

```bash
make cleanup-storage-orphans
```

This cleanup is used for files that are still `pending`, already expired, and have never been linked through `storage_file_links`.

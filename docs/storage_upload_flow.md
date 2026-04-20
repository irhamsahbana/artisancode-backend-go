# Storage Upload Flow

## Overview

Storage saat ini mendukung dua jalur upload:

1. presigned upload lewat `POST /storage/upload-url`
2. direct multipart upload lewat `POST /storage/upload`

Untuk attendance selfie dan flow mobile, jalur yang direkomendasikan tetap presigned upload karena backend menyimpan lifecycle file sebelum file di-attach ke resource bisnis.

## Main Tables

### `storage_files`

Menyimpan metadata file dan lifecycle upload:

- tenant ownership
- creator/uploader
- object key
- original filename
- content type
- status file seperti `pending`, `attached`, `deleted`

### `storage_file_links`

Menyimpan relasi file ke resource bisnis:

- `storage_file_id`
- `resource_type`
- `resource_id`
- `field_name`
- `sort_order`

Contoh attendance selfie:

- `resource_type = attendance_log`
- `field_name = selfie`

## Recommended Presigned Flow

### 1. Client request upload URL

Request:

- `POST /storage/upload-url`

Body minimal:

- `filename`
- `content_type`
- `folder`
- `is_public`

Backend akan:

- membentuk final object key
- membuat row `storage_files` dengan status awal `pending`
- mengembalikan `file_id`, `object_key`, `upload_url`, `method`, dan `headers`

### 2. Client upload binary ke object storage

Client meng-upload file ke URL hasil presign.

Penting:

- gunakan method dan headers yang diberikan backend
- content type harus sama dengan yang dipakai saat presign

### 3. Client submit business request

Contoh attendance:

- `POST /attendance-logs/check-in`
- `POST /attendance-logs/check-out`

Body attendance saat ini wajib menyertakan:

- `selfie_file_id`

### 4. Business record dibuat

Contoh:

- backend membuat `attendance_logs`

### 5. File di-link ke resource

Backend membuat row di `storage_file_links`.

### 6. File ditandai attached

Setelah link berhasil, backend akan memanggil `MarkFileAttached(...)` dan status file berpindah dari `pending` ke `attached`.

## Direct Multipart Flow

Route yang tersedia:

- `POST /storage/upload`

Flow ini langsung mendorong file ke storage integration tanpa lifecycle `presign -> attach` yang dipakai attendance.

Biasanya cocok untuk use case admin/internal yang hanya butuh upload file sederhana. Jika use case membutuhkan pengaitan file ke resource bisnis secara aman dan bisa dibersihkan bila orphan, lebih aman memakai flow presigned.

## Private File Access

File private diakses lewat:

- `GET /storage/private/*`

Route ini memakai signed URL validation middleware, lalu backend memetakan wildcard path menjadi object key private sebenarnya sebelum membuat file URL.

## Delete Behavior

Delete file memakai:

- `DELETE /storage/*`

Saat ini handler delete meneruskan object key/path ke core. File hanya boleh dihapus bila lolos aturan ownership dan state yang divalidasi core/repository.

## Orphan Cleanup

Cron task yang aktif:

```bash
go run ./cmd/bin/main.go cronjob --task=cleanup-expired-storage-files --limit=100
```

Atau via Make:

```bash
make cleanup-storage-orphans
```

Cleanup ini dipakai untuk file yang masih `pending`, sudah expired, dan belum pernah terhubung ke `storage_file_links`.

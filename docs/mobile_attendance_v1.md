# Mobile Attendance V1

Dokumen ini merangkum kontrak backend yang memang sudah tersedia sekarang untuk kebutuhan mobile attendance pertama.

## Goals

- login employee dengan auth flow yang sama
- check-in dan check-out mandiri oleh employee
- lihat histori attendance milik sendiri
- tampilkan konteks user, profil employee, shift hari ini, ringkasan hari ini, dan policy attendance

## Active Routes

Route yang sudah aktif saat ini:

- `GET /attendance-logs`
- `GET /attendance-logs/:id`
- `POST /attendance-logs/check-in`
- `POST /attendance-logs/check-out`
- `GET /attendance-summary/today`
- `GET /attendance-summary/owner-dashboard`
- `GET /attendance-policy`
- `GET /me`
- `GET /me/employee`
- `GET /me/shift-today`

Semua route di atas diproteksi `middleware.Auth`.

## Access Rules

### Employee

- hanya boleh melihat attendance miliknya sendiri
- hanya boleh check-in/check-out untuk dirinya sendiri
- `employee_id` ditentukan server dari `UserContext`, bukan dari request mobile

### Owner/Admin

- bisa melihat attendance tenant sesuai aturan core
- owner dashboard tersedia di `GET /attendance-summary/owner-dashboard`

## Attendance List Contract

### `GET /attendance-logs`

Query yang aktif sekarang:

- `page`
- `limit`
- `q`
- `employee_id`
- `type`
- `source`
- `status`
- `selfie_status`
- `org_unit_id`
- `branch_id`
- `work_location_id`
- `exception_type`
- `attendance_date`
- `date_from`
- `date_to`

Behavior umum:

- owner/admin bisa melihat data tenant sesuai akses
- employee dibatasi oleh core agar hanya melihat data sendiri

## Attendance Action Contract

### `POST /attendance-logs/check-in`

### `POST /attendance-logs/check-out`

Body request yang aktif sekarang:

```json
{
  "logged_at": "2026-03-28T08:05:00+08:00",
  "latitude": -5.1477,
  "longitude": 119.4327,
  "address": "Makassar",
  "device_id": "ios-sim-001",
  "device_name": "iPhone Test",
  "notes": "Arrived",
  "selfie_file_id": "uuidv7"
}
```

Catatan penting:

- `selfie_file_id` saat ini wajib
- `logged_at` optional, tapi kalau dikirim harus format RFC3339
- response sukses mengembalikan `id`

Success message saat ini:

- `"Attendance recorded successfully"`

## Supporting Mobile Routes

### `GET /me`

Dipakai untuk bootstrap identity session saat app start.

Field aktif:

- `user_id`
- `user_name`
- `tenant_id`
- `tenant_name`
- `roles`
- `company_id`
- `company_name`

### `GET /me/employee`

Dipakai untuk profil employee milik user yang login.

Field aktif:

- `id`
- `employee_no`
- `full_name`
- `email`
- `user_id`
- `org_unit_id`
- `job_position_id`
- `location_id`
- `shift_id`
- `status`
- `join_date`

### `GET /me/shift-today`

Dipakai untuk kartu shift hari ini.

Field aktif:

- `shift_id`
- `shift_name`
- `start_time`
- `end_time`
- `timezone`
- `attendance_date`

Jika user belum punya shift hari ini, endpoint saat ini bisa mengembalikan sukses dengan `data = null`.

### `GET /attendance-summary/today`

Dipakai untuk home summary mobile.

Field aktif:

- `attendance_date`
- `today_status`
- `checked_in`
- `checked_out`
- `check_in_log_id`
- `check_out_log_id`
- `last_log_type`
- `last_logged_at`
- `can_check_in`
- `can_check_out`

### `GET /attendance-policy`

Dipakai untuk menampilkan policy attendance yang relevan ke employee.

Field aktif:

- `timezone`
- `attendance_radius_meters`
- `attendance_check_in_start`
- `attendance_check_in_end`
- `attendance_check_out_start`
- `attendance_check_out_end`

## Storage Dependency

Flow mobile attendance saat ini bergantung pada storage module:

1. client minta upload URL ke `/storage/upload-url`
2. client upload selfie
3. client kirim `selfie_file_id` ke endpoint check-in/check-out

Jika `selfie_file_id` tidak valid, folder salah, atau status file bukan `pending`, request attendance akan ditolak.

## Business Rules Active Today

Dari implementasi core attendance saat ini:

1. satu employee hanya boleh punya satu `check_in` per tanggal attendance
2. satu employee hanya boleh punya satu `check_out` per tanggal attendance
3. `check_out` ditolak bila belum ada `check_in` pada tanggal yang sama
4. employee harus punya `shift_id`
5. selfie file harus ada, berasal dari folder attendance face, dan masih `pending`
6. source attendance mobile saat ini disimpan sebagai `mobile`
7. status attendance yang dibuat saat ini adalah `recorded`

## Known Scope Notes

- owner dashboard sudah tersedia, jadi web audit summary tidak lagi hanya backlog
- policy radius/window sudah tersedia via endpoint dedicated
- check-in/check-out masih online flow; belum ada offline sync atau idempotency layer khusus mobile

# AGENTS

Panduan backend dipisah ke folder `docs/`.

Kalau sebuah task mengubah workflow backend, konvensi coding, atau perilaku agent, update `agents.md` ini dan dokumen `docs/` yang relevan dalam task yang sama bila memungkinkan.

## Documentation Index

- [README](./docs/README.md)
- [Architecture](./docs/architecture.md)
- [Data Layers](./docs/data_layers.md)
- [DB & Migration](./docs/db_migration.md)
- [HTTP & Validation](./docs/http_validation.md)
- [Auth Context](./docs/auth_context.md)
- [Localization](./docs/localization.md)
- [Module Integration](./docs/module_integration.md)
- [Error Handling](./docs/error_handling.md)
- [Handler Pattern](./docs/handler_pattern.md)
- [Parameter Convention](./docs/parameter_convention.md)
- [Storage Upload Flow](./docs/storage_upload_flow.md)
- [Mobile Attendance V1](./docs/mobile_attendance_v1.md)
- [Auth Email Flow](./docs/auth_email_flow.md)

## Backend Reality Check

Arsitektur aktif backend saat ini:

- HTTP bootstrap ada di `internal/framework/primary/http/`
- HTTP dependency wiring ada di `internal/setup/http_dependency.go`
- Consumer bootstrap ada di `internal/framework/primary/consumer/postgres/`
- Consumer dependency wiring ada di `internal/setup/consumer_dependency.go`
- Business logic aktif ada di `internal/core/<module>`
- HTTP handler aktif ada di `internal/framework/primary/http/<module>`
- Postgres repository aktif ada di `internal/framework/secondary/db/postgres/<module>`
- Shared contracts tetap di `internal/ports/...`
- `internal/module/` saat ini bukan struktur dominan produksi; yang tersisa hanya template/eksperimen

Jangan mendokumentasikan `internal/module/<module>` sebagai pola utama backend kalau perubahanmu menyentuh kode produksi yang aktif sekarang.

## Quick Rules

### Layering

- Handler hanya menangani HTTP concern, validasi request, mapping, dan response.
- Core berisi business rule dan orkestrasi lintas dependency.
- Repository Postgres menangani query SQL dan mapping scan result.
- Handler tidak boleh mengirim `restentity` langsung ke core.
- Core tidak boleh mengimpor `restentity`.

### File Split Pattern

- Gunakan satu file per operasi utama di handler/core/repository bila modul memang sudah mengikuti pola split.
- `handler.go`, `core.go`, dan `repo.go` dipakai untuk struct, config, constructor, dan method registrasi dasar.
- Helper yang hanya dipakai oleh satu operasi boleh tinggal di file operasi itu.
- Helper shared lintas operasi boleh punya file helper sendiri bila itu membuat modul lebih jelas.

Contoh pola yang aktif sekarang:

```text
internal/framework/primary/http/attendance/
  handler.go
  get_attendance_logs.go
  get_attendance_log.go
  check_in.go
  check_out.go
  get_attendance_summary_today.go
  get_attendance_policy.go
  get_owner_attendance_dashboard.go
```

### Dependency Injection

- HTTP module wiring terpusat di `internal/setup/http_dependency.go`.
- Consumer wiring terpusat di `internal/setup/consumer_dependency.go`.
- Constructor gunakan config struct pattern.
- Dependency lintas modul harus lewat contract di `internal/ports/...`, bukan instantiate diam-diam di core/handler.

### Request Context

- Tenant scope ambil dari `common.GetUserContext(ctx)`, bukan dari request body/query.
- `restentity` tidak menyimpan `TenantID`.
- `common.UserContext` saat ini memuat:
  - `UserID`
  - `UserName`
  - `TenantID`
  - `TenantName`
  - `Roles`
  - `CompanyID`
  - `CompanyName`

### Tracing

- Semua function yang menerima `ctx` harus punya `tracing.StartSpan`.
- Format span:
  - handler HTTP: `internal:framework:primary:http:<module>:<file>:<function>`
  - core: `internal:core:<module>:<file>:<function>`
  - postgres repository: `internal:framework:secondary:db:postgres:<module>:<file>:<function>`

### Logging

- Pakai `log.Ctx(ctx)` kalau ada context.
- Parse/validation error di handler pakai `Warn`.
- Business rejection yang memang expected di core umumnya `Warn`.
- Error tak terduga atau kegagalan DB/integration pakai `Error`.

### Localization

- Jangan parse `Accept-Language` sendiri di feature module.
- HTTP app sudah memasang `WithRequestLanguage()` lalu `LocalizeJSONResponse()`.
- Gunakan `pkg/errmsg` supaya `message` dan `errors` ikut diterjemahkan otomatis.

### Current Public/Protected Bootstrap Notes

- `internal/framework/primary/http/build.go` mendaftarkan middleware global, metrics, dan memanggil `setup.HttpDependencies()`.
- User auth public routes diregister di `/users`.
- Sebagian besar route bisnis diproteksi dari `app.Group(..., middleware.Auth)` di setup layer.
- Storage punya kombinasi route protected dan route signed-public-ish untuk file private:
  - protected: `/storage/upload-url`, `/storage/upload`, `/storage/*`
  - signed URL validation: `/storage/private/*`

### Consumers and Jobs

- Consumer dependency aktif saat ini mencakup email subscription dan export job subscription.
- Cron task yang sudah ada mencakup:
  - `cleanup-expired-storage-files`
  - `process-export-jobs`
  - `cleanup-processed-message-queue`

Kalau menambah consumer baru atau cron task baru, dokumentasikan di `docs/module_integration.md` atau dokumen domain terkait.

# Architecture

## Active Runtime Structure

Struktur backend produksi yang aktif sekarang berpusat di layer berikut:

- `internal/framework/primary/http/`
  - bootstrap Fiber app
  - middleware global
  - HTTP handlers per module
- `internal/core/`
  - business logic per module
- `internal/framework/secondary/db/postgres/`
  - repository Postgres per module
- `internal/integration/`
  - adapter ke service eksternal seperti storage, email, token cache, OAuth, rate limit
- `internal/setup/`
  - dependency wiring untuk HTTP dan consumer runtime
- `internal/ports/`
  - contract antar layer

`internal/module/` bukan lagi struktur dominan kode produksi. Saat ini folder itu hanya menyisakan template/eksperimen, jadi jangan jadikan folder itu sebagai acuan utama dokumentasi arsitektur.

## Bootstrap Flow

### HTTP app

HTTP runtime dibangun dari:

- `internal/framework/primary/http/app.go`
- `internal/framework/primary/http/build.go`
- `internal/setup/http_dependency.go`

Urutan umumnya:

1. buat `fiber.App`
2. sinkronkan adapter global lewat `adapter.Adapters.Sync(...)`
3. pasang middleware global
4. expose `/metrics`
5. panggil `setup.HttpDependencies()` untuk me-register semua module route

### Consumer app

Consumer runtime saat ini dibangun dari:

- `internal/framework/primary/consumer/postgres/build.go`
- `internal/setup/consumer_dependency.go`

Runtime ini menyiapkan:

- message publisher
- subscription manager
- email subscription
- export job subscription
- export job core dan publisher dependency

## Module Shape

Pola modul aktif mengikuti pemisahan per concern, bukan satu folder `module/<name>`:

```text
internal/core/<module>
internal/framework/primary/http/<module>
internal/framework/secondary/db/postgres/<module>
internal/ports/core
internal/ports/secondary/db
```

Contoh konkret:

```text
internal/core/attendance/
internal/framework/primary/http/attendance/
internal/framework/secondary/db/postgres/attendance/
```

## File Split Pattern

Modul yang aktif umumnya memakai satu file per operasi utama.

Contoh:

```text
internal/core/user/
  core.go
  register.go
  login.go
  refresh_token.go
  verify_email.go
  resend_verification_email.go
  forgot_password.go
  reset_password.go
  action_token_helpers.go
```

Aturan praktis:

- `handler.go`, `core.go`, `repo.go` dipakai untuk struct, config, constructor, dan method registrasi
- satu operasi utama idealnya satu file
- helper yang spesifik ke satu operasi boleh diletakkan di file yang sama
- helper shared boleh punya file sendiri bila dipakai lintas operasi

## Ports

Interface tetap disentralisasi di `internal/ports/`:

- `internal/ports/core/`
- `internal/ports/primary/`
- `internal/ports/secondary/db/`
- `internal/ports/secondary/integration/`

Gunakan contract ini untuk dependency injection antar layer dan antar modul.

## Mapper

Mapping antar boundary tetap dipisah di `internal/entity/mapper/`.

Peran mapper:

- `restentity` -> `coreentity`
- `repoentity` -> `coreentity`
- `coreentity` -> `restentity`

Handler sebaiknya mapping di boundary, bukan membiarkan core menerima `restentity` langsung.

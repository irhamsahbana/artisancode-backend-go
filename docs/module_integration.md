# Module Integration

## Where Wiring Lives Now

Module registration tidak lagi cukup dijelaskan sebagai satu file `http_dependency.go`.

Wiring aktif sekarang terbagi dua:

- `internal/setup/http_dependency.go`
  - repository, core, handler, dan route registration untuk HTTP runtime
- `internal/setup/consumer_dependency.go`
  - subscription, consumer-facing core, publisher, dan shutdown wiring untuk worker/consumer runtime

Kalau menambah module HTTP baru, update `http_dependency.go`.
Kalau menambah subscription/consumer baru, update `consumer_dependency.go`.

## HTTP Registration Pattern

`internal/framework/primary/http/build.go` hanya menyiapkan app dan middleware global, lalu memanggil `setup.HttpDependencies()`.

Di `setup.HttpDependencies()` pola registration yang dipakai sekarang adalah:

1. build repository
2. build integration/helper dependency
3. build core
4. build handler
5. register route dengan `app.Group(...)`

Contoh bentuk yang dipakai:

```go
employeeRepository := employeeRepo.NewEmployeeRepository(employeeRepo.EmployeeRepositoryConfig{
    DB: db,
})

employeeCoreInst := employeeCore.NewEmployeeCore(employeeCore.EmployeeCoreConfig{
    Repo:     employeeRepository,
    UserRepo: userRepository,
})

employeeHandler.NewEmployeeHandler(employeeHandler.EmployeeHandlerConfig{
    Core: employeeCoreInst,
}).Register(app.Group("/employees", middleware.Auth))
```

## Route Protection

Proteksi route saat ini kebanyakan dipasang di setup layer lewat `app.Group(path, middleware.Auth)`, bukan di dalam setiap handler.

Contoh yang aktif:

- `app.Group("/companies", middleware.Auth)`
- `app.Group("/employees", middleware.Auth)`
- `app.Group("/attendance-logs", middleware.Auth)`
- `app.Group("/me", middleware.Auth)`

Pengecualian yang memang dibiarkan public atau mixed:

- `/users/*`
  - auth/public endpoints diregister di handler user
  - CRUD tertentu diproteksi per-route di handler
- `/storage`
  - upload/delete/list diproteksi
  - `GET /storage/private/*` memakai signed URL validation middleware

Karena itu, saat menambah route baru:

- kalau seluruh group harus protected, pasang `middleware.Auth` di `app.Group(...)`
- kalau satu module punya kombinasi public dan protected route, dokumentasikan dengan jelas di `handler.go`

## Consumer Registration Pattern

Consumer runtime saat ini tidak me-register HTTP route. Yang dibuat adalah dependency untuk subscription manager dan handler loop.

`NewConsumerDependencies(...)` saat ini menyiapkan:

- email subscription
- export job subscription
- export job core
- export job publisher
- subscription manager
- shutdown callback

Jika menambah consumer baru:

1. buat subscription config baru
2. inject dependency yang dibutuhkan core/processor
3. expose hasilnya di `ConsumerDependencies`
4. assign hasilnya di `internal/framework/primary/consumer/postgres/build.go`

## Cross-Module Dependencies

Gunakan interface dari `internal/ports/...` untuk dependency lintas modul.

Contoh yang memang aktif:

- `employee` core memakai `UserRepository`
- `userinvitation` core memakai `UserRepository` dan `EmployeeRepository`
- `worklocation` core memakai `OrgUnitRepository`
- `attendance` core memakai `CompanyRepository` dan `StorageRepository`
- `export_job` core memakai `AttendanceRepository`, `StorageRepository`, dan message publisher

Catatan invitation flow:

- module `user-invitations` diregister di `internal/setup/http_dependency.go`
- route `GET /user-invitations/accept` dan `POST /user-invitations/accept` bersifat public
- route create/list/resend/revoke invitation diproteksi auth di handler per-route

Catatan employee flow:

- `CreateEmployee` tidak lagi otomatis membuat record `users`
- akses login employee sekarang diharapkan lewat invitation/activation flow terpisah

Aturan:

- jangan instantiate repository atau integration baru di dalam handler/core
- jangan akses package concrete module lain langsung kalau contract sudah tersedia di `ports`
- semua wiring lintas modul harus terlihat di `setup`

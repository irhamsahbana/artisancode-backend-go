# Module Integration

## Where Wiring Lives Now

Module registration can no longer be explained as a single `http_dependency.go` file.

Active wiring is now split into two places:

- `internal/setup/http_dependency.go`
  - repository, core, handler, and route registration for the HTTP runtime
- `internal/setup/consumer_dependency.go`
  - subscription, consumer-facing core, publisher, and shutdown wiring for the worker/consumer runtime

If you add a new HTTP module, update `http_dependency.go`.
If you add a new subscription or consumer, update `consumer_dependency.go`.

## HTTP Registration Pattern

`internal/framework/primary/http/build.go` only prepares the app and global middleware, then calls `setup.HttpDependencies()`.

Inside `setup.HttpDependencies()`, the current registration pattern is:

1. build repository
2. build integration/helper dependency
3. build core
4. build handler
5. register routes with `app.Group(...)`

Example:

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

Route protection is currently applied mostly in the setup layer through `app.Group(path, middleware.Auth)`, not inside every handler.

Current examples:

- `app.Group("/companies", middleware.Auth)`
- `app.Group("/employees", middleware.Auth)`
- `app.Group("/attendance-logs", middleware.Auth)`
- `app.Group("/me", middleware.Auth)`

Known exceptions that intentionally stay public or mixed:

- `/users/*`
  - auth/public endpoints are registered in the user handler
  - some CRUD routes are protected per-route in the handler
- `/internal-commerce/*`
  - protected by `middleware.InternalAuth`
  - internal-only commerce, quotation, order, invoice, and payment operations
- `/webhooks/doku`
  - public webhook endpoint for DOKU payment callbacks
- `/storage`
  - upload, delete, and list routes are protected
  - `GET /storage/private/*` uses signed URL validation middleware

Because of that, when adding a new route:

- if the entire group must be protected, install `middleware.Auth` in `app.Group(...)`
- if one module has a mix of public and protected routes, document that clearly in `handler.go`

## Consumer Registration Pattern

The consumer runtime does not register HTTP routes. It builds dependencies for the Watermill router and registers one consumer handler per topic.

The primary Postgres consumer structure is now split per feature to stay consistent with the HTTP split pattern:

- `internal/framework/primary/consumer/postgres/`
  - consumer runtime bootstrap (`app.go`, `build.go`, `run.go`, `shutdown.go`)
- `internal/framework/primary/consumer/postgres/user/`
  - consumer for email verification and forgot password
- `internal/framework/primary/consumer/postgres/userinvitation/`
  - consumer for invitation email
- `internal/framework/primary/consumer/postgres/export_job/`
  - consumer for export jobs
- `internal/framework/primary/consumer/postgres/shared/`
  - shared helpers for message carrier headers and consumer email/logging utilities

`NewConsumerDependencies(...)` currently prepares:

- export job core
- shutdown callback

Message subscriptions use Watermill SQL default Postgres Pub/Sub schema. Each topic has its own `watermill_<topic>` table and `watermill_offsets_<topic>` table, with production schema managed manually by goose migrations. Subscribers must set an explicit consumer group such as `notification_service`, `stock_service`, or `audit_service`.

The current consumer router also installs shared middleware for:

- panic recovery
- bounded retry using message-bus retry config
- poison queue publish to the shared `dead_letter` topic after retries are exhausted

When adding a new consumer:

1. create topic constants in `internal/entity/common/message_bus.go`
2. inject the dependencies required by the core/processor
3. expose the result through `ConsumerDependencies` when the handler needs new core dependencies
4. register a dedicated Watermill router consumer handler in `internal/framework/primary/consumer/postgres/build.go`
5. place handlers in the matching feature consumer folder instead of dropping all handler files into the root `consumer/postgres` folder

## Scheduler Registration Pattern

The scheduler runtime is separate from both HTTP and consumer bootstraps. `cmd/scheduler.go` builds the gocron runtime, wires repository dependencies, and registers long-running jobs.

Current scheduler jobs include:

- storage cleanup
- message queue cleanup
- billing renewal processing
- billing dunning escalation

Billing scheduler logic currently lives under `internal/scheduler/billing/`.

When adding a new scheduler job:

1. register the job in `cmd/scheduler.go`
2. keep domain execution logic under `internal/scheduler/<domain>/`
3. drive cadence and batch or grace limits from config/env instead of hardcoding them in the job body

## Cross-Module Dependencies

Use interfaces from `internal/ports/...` for cross-module dependencies.

## Integration Adapter Pattern

For outbound provider clients and other third-party service adapters, use:

- `internal/integration/<provider>`
  - concrete client, request signing, response parsing, webhook verification, and transport helpers

If an integration is used by core or handlers:

1. instantiate it in `internal/setup/http_dependency.go` or the relevant runtime setup
2. inject it into core through the relevant contract/interface
3. do not build third-party HTTP clients directly inside handlers or core

Example:

```text
internal/integration/doku/
  client.go
  client_test.go
```

Current examples:

- `employee` core uses `UserRepository`
- `userinvitation` core uses `UserRepository` and `EmployeeRepository`
- `worklocation` core uses `OrgUnitRepository`
- `attendance` core uses `CompanyRepository` and `StorageRepository`
- `export_job` core uses `AttendanceRepository`, `StorageRepository`, and a message publisher

Invitation flow notes:

- the `user-invitations` module is registered in `internal/setup/http_dependency.go`
- `GET /user-invitations/accept` and `POST /user-invitations/accept` are public routes
- create, list, resend, and revoke invitation routes are auth-protected per-route in the handler

Employee flow notes:

- `CreateEmployee` no longer creates a `users` record automatically
- employee login access is now expected to go through a separate invitation/activation flow

Rules:

- do not instantiate new repositories or integrations inside handlers or core
- do not access another module's concrete package directly when a contract already exists in `ports`
- all cross-module wiring must remain visible in `setup`

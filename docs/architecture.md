# Architecture

## Active Runtime Structure

The active production backend currently centers on these layers:

- `internal/framework/primary/http/`
  - Fiber app bootstrap
  - global middleware
  - HTTP handlers per module
- `internal/core/`
  - business logic per module
- `internal/framework/secondary/db/postgres/`
  - Postgres repositories per module
- `internal/integration/`
  - adapters for external services such as storage, email, token cache, OAuth, rate limiting, and third-party provider clients like DOKU
- `internal/setup/`
  - dependency wiring for HTTP and consumer runtimes
- `internal/ports/`
  - contracts between layers

`internal/module/` is no longer the dominant production structure. It now mostly contains templates and experiments, so do not use it as the main reference for architecture documentation.

## Bootstrap Flow

### HTTP app

The HTTP runtime is built from:

- `internal/framework/primary/http/app.go`
- `internal/framework/primary/http/build.go`
- `internal/setup/http_dependency.go`

Typical sequence:

1. create `fiber.App`
2. sync global adapters through `adapter.Adapters.Sync(...)`
3. install global middleware
4. expose `/metrics`
5. call `setup.HttpDependencies()` to register all module routes

### Consumer app

The current consumer runtime is built from:

- `internal/framework/primary/consumer/postgres/build.go`
- `internal/setup/consumer_dependency.go`

This runtime currently prepares:

- message publisher
- subscription manager
- email subscription
- export job subscription
- export job core and publisher dependencies

## Module Shape

Active modules follow separation by concern, not a single `module/<name>` folder:

```text
internal/core/<module>
internal/framework/primary/http/<module>
internal/framework/secondary/db/postgres/<module>
internal/ports/core
internal/ports/secondary/db
```

Concrete example:

```text
internal/core/attendance/
internal/framework/primary/http/attendance/
internal/framework/secondary/db/postgres/attendance/
```

## File Split Pattern

Active modules generally use one file per main operation.

Example:

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

Practical rules:

- `handler.go`, `core.go`, and `repo.go` are for struct definitions, config, constructors, and registration methods
- one main operation should ideally live in one file
- helpers specific to one operation may stay in the same file
- shared helpers may live in their own file when reused across operations

## Ports

Interfaces remain centralized in `internal/ports/`:

- `internal/ports/core/`
- `internal/ports/secondary/db/`
- `internal/ports/integration/`

Use these contracts for dependency injection across layers and across modules.

## Mapper

Mapping across boundaries remains separated in `internal/entity/mapper/`.

Mapper responsibilities:

- `restentity` -> `coreentity`
- `repoentity` -> `coreentity`
- `coreentity` -> `restentity`

Handlers should perform mapping at the boundary instead of allowing core to receive `restentity` directly.

## External Integrations

Outbound provider adapters live in `internal/integration/<provider>`.

Example:

```text
internal/integration/doku/
  client.go
  client_test.go
```

Keep provider-specific request signing, response parsing, webhook signature verification, and transport concerns inside the integration package rather than spreading them across handlers or core.

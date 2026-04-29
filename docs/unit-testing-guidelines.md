# Backend Unit Testing Guidelines

## 1. Required Tools

- Use `github.com/stretchr/testify/require` for assertions.
- Use `mockery` as the primary mock generator.
- Use generated testify expecter mocks through `mock.EXPECT()`.
- Do not use `gomock`.
- Do not hand-write mocks unless an interface cannot reasonably be generated.

## 2. Mock Generation

Mocks are configured in the backend root `.mockery.yaml` using mockery's modern `packages` format. The Makefile pins mockery to `v3.7.0` so generation does not depend on whichever `mockery` binary appears first in the local `PATH`.

Generate mocks with:

```bash
make mock
```

Mockery v3's `template: testify` generates expecter-style mocks. The generated mocks expose this style:

```go
repo := mocks.NewUserRepository(t)

repo.EXPECT().
    FindByID(ctx, "user-1").
    Return(&entity.User{ID: "user-1", Name: "Budi"}, nil)
```

Place generated mocks under the related port package's `mocks` subfolder, for example:

```text
internal/ports/secondary/db/mocks
internal/ports/integration/mocks
```

This keeps mocks out of production packages and avoids import cycles.

## 3. What To Mock

Mock dependency interfaces used by usecase/service/application logic:

- Repository interfaces
- Publisher interfaces
- External API client interfaces
- Transaction manager interfaces
- Cache interfaces
- Feature flag interfaces

Do not generate mocks for:

- Entity/model structs
- DTO/request/response structs
- Helper functions
- Concrete repository implementations
- Direct database clients
- Standard library packages

If an interface is too large and forces broad, fragile expectations, recommend splitting it into smaller role-based interfaces. Do not perform a large interface refactor unless the task explicitly needs it.

## 4. Test Scope

Prioritize unit tests for active business logic in:

```text
internal/core/<module>
```

Keep HTTP handler tests focused on HTTP behavior, request validation, mapping, and response shape. Keep repository tests separate from usecase tests; repository tests should not depend on generated usecase mocks.

## 5. Test File Structure

Place test files beside the code under test and use `_test.go` filenames.

Example:

```text
internal/core/user/create_user.go
internal/core/user/create_user_test.go
```

## 6. Test Style

Use table-driven tests when there is more than one case. Always use `t.Run` for named cases.

Use `require` assertions:

```go
require.NoError(t, err)
require.Error(t, err)
require.Equal(t, expected, actual)
```

Do not over-mock internal implementation details. Set expectations only for dependencies that are part of the observable usecase path for that case.

## 7. Minimum Coverage Per Usecase Test

For a usecase/service test, cover at least:

- Success case
- Validation/business rejection case
- Dependency error case
- Important edge case specific to the business logic

## 8. Verification

Run the touched package first, then the full backend suite before handing off:

```bash
go test ./internal/core/<module>
go test ./...
```

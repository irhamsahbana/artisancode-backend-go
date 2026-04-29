# Backend Unit Testing Guidelines

This guide defines the backend unit testing standard. The main goal of unit testing is to make small business behaviors trustworthy, fast to run, deterministic, and easy to understand when they fail.

## 1. General Unit Testing Principles

Good unit tests should be:

- **Fast**: avoid real databases, networks, external services, or heavy file system usage.
- **Deterministic**: the result must always be the same and must not depend on current time, random values, test order, machine timezone, or local environment.
- **Focused**: each test case should ideally verify one main behavior.
- **Readable**: test names and case names should explain the business scenario.
- **Isolated**: one test must not depend on another test.
- **Behavior-oriented**: assert observable output, errors, and side effects instead of internal implementation details.

Do not write tests only to increase coverage. Prioritize tests that protect important business behavior.

## 2. Required Tools

- Use `github.com/stretchr/testify/require` for assertions.
- Use `mockery` as the primary mock generator.
- Use generated testify expecter mocks through `mock.EXPECT()`.
- Do not use `gomock`.
- Do not hand-write mocks unless an interface cannot reasonably be generated.

## 3. Mock Generation

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

## 4. What To Mock

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
- Pure functions or value objects

If a helper is deterministic and cheap, call the real implementation instead of mocking it.

If an interface is too large and forces broad, fragile expectations, recommend splitting it into smaller role-based interfaces. Do not perform a large interface refactor unless the task explicitly needs it.

## 5. Test Doubles

Use the simplest test double that proves the behavior:

- **Stub**: returns fixed data.
- **Mock**: verifies dependency interaction.
- **Fake**: lightweight working implementation used only for tests.
- **Spy**: records calls for later assertion.

Prefer mocks when interaction matters.
Prefer fakes or stubs when only returned data matters and strict expectations would make the test brittle.

## 6. Test Scope

Prioritize unit tests for active business logic in:

```text
internal/core/<module>
```

Keep HTTP handler tests focused on:

- HTTP status code
- Request validation
- Request mapping
- Response shape
- Auth/middleware behavior when relevant

Keep repository tests separate from usecase tests. Repository tests should not depend on generated usecase mocks.

## 7. Test File Structure

Place test files beside the code under test and use `_test.go` filenames.

Example:

```text
internal/core/user/create_user.go
internal/core/user/create_user_test.go
```

## 8. Test Style

Use table-driven tests when there is more than one similar case. Always use `t.Run` for named cases.

Use `require` assertions:

```go
require.NoError(t, err)
require.Error(t, err)
require.Equal(t, expected, actual)
```

Do not over-mock internal implementation details. Set expectations only for dependencies that are part of the observable usecase path for that case.

## 9. Recommended Test Structure

Structure usecase/service tests as Arrange, Act, Assert:

- **Arrange**: create the context, inputs, generated mocks, mock expectations, and the core/service under test.
- **Act**: call exactly one public method on the core/service.
- **Assert**: check the error, output, and observable side effects with `require`.

For table-driven tests, keep case-specific mock behavior inside a `setup` function:

```go
func TestUsecaseName(t *testing.T) {
    ctx := context.Background()

    validInput := coreentity.SomeInput{ID: "item-1"}

    tests := []struct {
        name      string
        input     coreentity.SomeInput
        setup     func(repo *mocks.SomeRepository, bus *mocks.MessagePublisher)
        want      *coreentity.SomeResult
        wantError bool
    }{
        {
            name:  "success",
            input: validInput,
            setup: func(repo *mocks.SomeRepository, bus *mocks.MessagePublisher) {
                repo.EXPECT().
                    FindByID(mock.Anything, "item-1").
                    Return(&coreentity.SomeResult{ID: "item-1"}, nil)
            },
            want: &coreentity.SomeResult{ID: "item-1"},
        },
        {
            name:      "validation error when id is empty",
            input:     coreentity.SomeInput{ID: ""},
            wantError: true,
        },
        {
            name:  "dependency error from repository",
            input: validInput,
            setup: func(repo *mocks.SomeRepository, bus *mocks.MessagePublisher) {
                repo.EXPECT().
                    FindByID(mock.Anything, "item-1").
                    Return(nil, errors.New("repository failed"))
            },
            wantError: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            repo := mocks.NewSomeRepository(t)
            bus := mocks.NewMessagePublisher(t)

            if tt.setup != nil {
                tt.setup(repo, bus)
            }

            core := NewSomeCore(Config{
                Repo: repo,
                Bus:  bus,
            })

            got, err := core.DoSomething(ctx, tt.input)

            if tt.wantError {
                require.Error(t, err)
                return
            }

            require.NoError(t, err)
            require.Equal(t, tt.want, got)
        })
    }
}
```

Prefer `mock.Anything` for `context.Context` expectations because tracing can wrap the incoming context before it reaches dependencies. Keep exact matching for business arguments that matter to the behavior.

## 10. Minimum Coverage Per Usecase Test

For a usecase/service test, cover at least:

- Success case
- Validation/business rejection case
- Dependency error case
- Important edge case specific to the business logic

For mutation usecases, also consider:

- Duplicate/idempotency case
- Authorization or ownership rejection case
- Transaction rollback case
- Publish/event failure case
- Cache invalidation behavior

Do not force all cases if they are not relevant to the usecase.

## 11. Error Assertions

Avoid only checking whether an error exists when the specific error matters.

Prefer `require.ErrorIs` when the code returns sentinel/wrapped errors:

```go
require.ErrorIs(t, err, entity.ErrInvalidStatus)
```

Prefer `require.ErrorAs` when the code returns typed errors:

```go
var validationErr *entity.ValidationError
require.ErrorAs(t, err, &validationErr)
```

Use `wantError bool` only when the exact error type/value is not important.

Recommended table field when error identity matters:

```go
tests := []struct {
    name    string
    input   coreentity.SomeInput
    setup   func(repo *mocks.SomeRepository)
    want    *coreentity.SomeResult
    wantErr error
}{
    {
        name:    "invalid status",
        input:   coreentity.SomeInput{Status: "unknown"},
        wantErr: entity.ErrInvalidStatus,
    },
}
```

Assertion:

```go
if tt.wantErr != nil {
    require.Error(t, err)
    require.ErrorIs(t, err, tt.wantErr)
    return
}

require.NoError(t, err)
```

## 12. Deterministic Inputs

Unit tests must not depend on:

- `time.Now()` directly
- Random UUIDs
- Random numbers
- Network calls
- Real database state
- External service state
- Test execution order
- Machine-local timezone/configuration

If production code depends on time, UUIDs, random values, or external configuration, inject them or wrap them behind a small dependency.

Prefer:

- Fixed timestamps
- Fixed UUIDs
- `t.TempDir()` for temporary files
- Explicit timezone handling
- Stable test fixtures

Example:

```go
fixedNow := time.Date(2026, 4, 30, 10, 0, 0, 0, time.UTC)
```

If the code needs `now`, prefer injecting a clock:

```go
type Clock interface {
    Now() time.Time
}
```

## 13. Side Effects

When a usecase produces side effects, assert the observable side effect through dependency expectations.

Examples:

- Repository update is called with the expected status.
- Publisher publishes the expected event/topic/payload.
- Transaction manager wraps the expected mutation.
- Cache invalidation happens only after successful mutation.
- External API client is called only after validation passes.

Avoid asserting irrelevant internal call order unless order is part of the business rule.

Example:

```go
publisher.EXPECT().
    PublishOrderCreated(mock.Anything, mock.MatchedBy(func(event entity.OrderCreatedEvent) bool {
        return event.OrderID == "order-1" && event.CompanyID == "company-1"
    })).
    Return(nil)
```

## 14. Negative Expectations

For validation/business rejection cases, ensure downstream dependencies are not called when they should not be reached.

In testify/mock, unexpected calls usually fail automatically if no expectation is registered. Therefore, for a validation error case, prefer not setting any repository/publisher expectation unless the call is expected.

Example:

```go
{
    name:      "validation error should not call repository",
    input:     coreentity.SomeInput{ID: ""},
    setup:     nil,
    wantError: true,
}
```

If a method may or may not be called and the call is not relevant to the behavior, reconsider whether the test is too coupled to implementation details.

## 15. Test Data Builders

Use test data builders for large entities or request structs. Keep only the fields relevant to the test case visible in each test.

Example:

```go
func validInput() coreentity.SomeInput {
    return coreentity.SomeInput{
        ID:        "item-1",
        UserID:    "user-1",
        CompanyID: "company-1",
        Status:    "active",
    }
}
```

Then override only fields relevant to the scenario:

```go
input := validInput()
input.Status = "unknown"
```

This keeps tests readable and prevents every test case from duplicating large setup data.

## 16. Table-Driven Test Guidance

Use table-driven tests when cases share the same Arrange/Act/Assert shape.

Do not force table-driven tests when each case has very different setup or assertions. In that situation, separate explicit tests are easier to read and maintain.

Good use of table-driven tests:

- Multiple validation cases
- Multiple business rule branches with similar setup
- Same dependency pattern with different inputs/outputs

Avoid table-driven tests when:

- The case struct becomes too large
- Each `setup` function is completely different
- Assertions require many custom branches
- The test becomes harder to read than separate test functions

## 17. Parallel Tests

Use `t.Parallel()` only when the test does not mutate shared state.

Avoid `t.Parallel()` when the test touches:

- Package-level variables
- Global config
- Environment variables
- Shared fixtures
- Time mocks
- Random seed
- Shared database/container
- External service stub with shared state

If using `t.Setenv`, be careful with parallel tests because environment variables are process-wide.

## 18. Naming Guidelines

Use names that describe behavior and business scenario.

Prefer:

```go
func TestCreateOrder_ReturnsValidationError_WhenCustomerIDIsEmpty(t *testing.T)
```

or table case names like:

```go
name: "returns validation error when customer id is empty"
```

Avoid vague names:

```go
name: "case 1"
name: "failed"
name: "error"
```

## 19. What Not To Test In Unit Tests

Avoid testing these through unit tests:

- Third-party library internals
- Generated mock code
- Framework routing internals
- Database query correctness
- Real message broker behavior
- Real network behavior

Use integration tests for database queries, repository implementations, real broker behavior, and framework wiring.

## 20. Verification

Run the touched package first, then the full backend suite before handing off:

```bash
go test ./internal/core/<module>
go test ./...
```

When touching concurrent code:

```bash
go test -race ./...
```

When checking for hidden cache/order issues:

```bash
go test -count=1 ./...
go test -shuffle=on ./...
```

Use `-count=1` when you want to avoid cached results.
Use `-shuffle=on` occasionally to detect hidden test-order dependencies.
Use `-race` when touching goroutines, shared memory, worker pools, async publisher logic, or graceful shutdown behavior.

## 21. Review Checklist

Before considering a unit test done, check:

- Does the test protect an important behavior?
- Is the test deterministic?
- Does it avoid real network/database/external service calls?
- Is the setup readable?
- Are only relevant dependencies mocked?
- Does the test assert the important output/error/side effect?
- Does an error case verify the correct error when it matters?
- Does a validation rejection avoid unnecessary downstream calls?
- Would the test still pass after a safe internal refactor?
- Can another developer understand the failed test quickly?

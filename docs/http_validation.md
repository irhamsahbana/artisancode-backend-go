# HTTP & Validation

## Foreign Key Rules in Create/Update Requests

- Required FK: `validate:"required,uuidv7"`.
- Optional FK: `validate:"omitempty,uuidv7"`.

## Handler Pattern

- Parse request using `ParamsParser/QueryParser/BodyParser`.
- Validate using `adapter.Adapters.Validator`.
- Map errors using `pkg/errmsg`.
- Return responses using `pkg/response`.

For detailed handler templates, see [Handler Pattern](./handler_pattern.md).

## Request Struct Conventions

### Path Parameters

```go
type GetItemReq struct {
    ID string `params:"id" validate:"required"`
}
```

### Query Parameters (List)

```go
type GetItemsReq struct {
    Q string `query:"q" validate:"omitempty,min=2"`
    types.MetaQuery  // embeds Page, Paginate
}

func (r *GetItemsReq) SetDefault() {
    r.MetaQuery.SetDefault()
}
```

### Body (Create)

```go
type CreateItemReq struct {
    Name    string  `json:"name" validate:"required,min=2"`
    ParentID *string `json:"parent_id" validate:"omitempty,uuidv7"`
}
```

### Body + Params (Update)

```go
type UpdateItemReq struct {
    ID   string `params:"id" validate:"required"`
    Name string `json:"name" validate:"required,min=2"`
}
```

## Validation Tags Reference

| Tag | Usage |
|-----|-------|
| `required` | Field must be present |
| `omitempty` | Skip validation if empty |
| `uuidv7` | Must be valid UUIDv7 |
| `email` | Must be valid email format |
| `min=N` | Minimum length (string) or value (number) |
| `max=N` | Maximum length (string) or value (number) |
| `oneof=a b c` | Must be one of listed values |

## Auth Email Requests

Auth email endpoints currently use these request shapes:

```go
type VerifyEmailReq struct {
    Token string `json:"token" validate:"required"`
}

type ResendVerificationEmailReq struct {
    Email string `json:"email" validate:"required,email"`
}

type ForgotPasswordReq struct {
    Email string `json:"email" validate:"required,email"`
}

type ResetPasswordReq struct {
    Token    string `json:"token" validate:"required"`
    Password string `json:"password" validate:"required,min=8"`
}
```

Notes:

- resend verification and forgot password are public endpoints
- reset password and verify email consume one-time action tokens
- rate-limited handlers should return `429` with `Retry-After`

## Domain Constants

- When a request or response field maps to a known domain enum such as attendance `type`, `source`, or `status`, define and reuse the canonical constant in `internal/entity/common/enum.go`.
- Validation may still use `oneof=...` at the HTTP boundary, but handler/core/repository code should avoid repeating the same raw strings once a shared constant exists.

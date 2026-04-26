# Error Handling

## Error Types

The application uses three error handling mechanisms via `pkg/errmsg`:

### 1. Custom Errors (`errmsg.NewCustomErrors`)

For business logic and domain errors with specific HTTP status codes:

```go
// Not found
return nil, errmsg.NewCustomErrors(404).SetMessage("Employee not found")

// Bad request
return nil, errmsg.NewCustomErrors(400).SetMessage("Email is already registered")

// Unauthorized
return nil, errmsg.NewCustomErrors(401).SetMessage("Invalid or expired refresh token")

// Internal server error
return nil, errmsg.NewCustomErrors(500).SetMessage("Failed to generate token")

// With field-level errors
err := errmsg.NewCustomErrors(400).
    Add("file", "file is required").
    SetMessage("file is required")
```

### 2. Validation Errors

Automatically handled by `errmsg.Errors(ctx, err, req)` when `err` is `validator.ValidationErrors`:

```go
if err := v.Validate(req); err != nil {
    code, errors := errmsg.Errors(ctx, err, req)
    return c.Status(code).JSON(response.Error(errors))
}
```

### 3. Database Errors

Automatically handled by `errmsg.Errors[error](ctx, err)` when `err` is `*pgconn.PgError`.

## Status Code Guidelines

| Code | When to Use | Example |
|------|-------------|---------|
| 400 | Invalid input, business rule violation | Duplicate email, invalid credentials |
| 401 | Authentication failure | Expired token, invalid refresh token |
| 403 | Authorization failure | Insufficient role/permission |
| 404 | Resource not found | Employee/company not found |
| 500 | Unexpected server error | Token generation failure, DB connection error |

## Error Handling by Layer

### Repository Layer

- Return `errmsg.NewCustomErrors(404)` for `sql.ErrNoRows`
- Log `Warn` before returning expected custom errors from repository code, including `sql.ErrNoRows`, `RowsAffected() == 0`, and repository-level validation failures.
- Log `Error` before returning unexpected DB/integration errors, including failed queries, failed `ExecContext`, failed `RowsAffected()`, failed marshal/unmarshal caused by persisted data, and transaction begin/commit failures.
- Let other DB errors propagate naturally after logging (errmsg handles `*pgconn.PgError`)
- Do not warn-log benign existence checks or optional lookups that intentionally return `false, nil` or `nil, nil`.

```go
err := r.db.GetContext(ctx, &data, r.db.Rebind(query), filter.ID, filter.TenantID)
if err != nil {
    if err == sql.ErrNoRows {
        log.Ctx(ctx).Warn().Any(common.LogKeyPayload, filter).Msg("Employee not found")
        return nil, errmsg.NewCustomErrors(404).SetMessage("Employee not found")
    }
    log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to get employee")
    return nil, err
}
```

For update/delete operations, treat zero affected rows as an expected miss and log it as `Warn`:

```go
result, err := r.db.ExecContext(ctx, r.db.Rebind(query), filter.ID, filter.TenantID)
if err != nil {
    log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to delete employee")
    return err
}

rowsAffected, err := result.RowsAffected()
if err != nil {
    log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to get delete employee rows affected")
    return err
}
if rowsAffected == 0 {
    log.Ctx(ctx).Warn().Any(common.LogKeyPayload, filter).Msg("Employee not found when deleting")
    return errmsg.NewCustomErrors(404).SetMessage("Employee not found")
}
```

### Core Layer

- Return `errmsg.NewCustomErrors` for business rule violations
- Log expected business rejections with `Warn` before returning custom errors
- Propagate repository errors without wrapping

```go
exist, err := c.repo.ExistsActiveUserByEmail(ctx, user.Email)
if err != nil {
    return nil, err  // propagate repo error
}
if exist {
    log.Ctx(ctx).Warn().Any(common.LogKeyPayload, map[string]string{"email": user.Email}).Msg("Email already registered")
    return nil, errmsg.NewCustomErrors(400).SetMessage("Email is already registered")
}
```

### Handler Layer

- Use `errmsg.Errors(ctx, err, req)` for validation errors (with payload for field-level mapping)
- Use `errmsg.Errors[error](ctx, err)` for core/repo errors (without payload)

```go
// Validation error — pass ctx and req for language + field mapping
code, errors := errmsg.Errors(ctx, err, req)

// Core/repo error — ctx provides request language
code, errs := errmsg.Errors[error](ctx, err)
```

## Response Envelope

All responses follow this format:

```json
// Success
{
    "success": true,
    "message": "Your request has been successfully processed",
    "data": { ... }
}

// Error
{
    "success": false,
    "message": "Error message here",
    "errors": {
        "field_name": ["error message"]
    }
}
```

Use `response.Success(data, "")` and `response.Error(err)` from `pkg/response`.

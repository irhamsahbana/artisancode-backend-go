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

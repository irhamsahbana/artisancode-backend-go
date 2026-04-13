# Handler Pattern

## CRUD Handler Structure

Every CRUD module handler follows a standard structure.

### Constructor

```go
type employeeHandler struct {
    core portsCore.EmployeeCore
}

type EmployeeHandlerConfig struct {
    Core portsCore.EmployeeCore
}

func NewEmployeeHandler(cfg EmployeeHandlerConfig) *employeeHandler {
    return &employeeHandler{core: cfg.Core}
}
```

### Route Registration

```go
func (h *employeeHandler) Register(router fiber.Router) {
    router.Get("/", h.getEmployees)
    router.Get("/:id", h.getEmployee)
    router.Post("/", h.createEmployee)
    router.Put("/:id", h.updateEmployee)
    router.Delete("/:id", h.deleteEmployee)
}
```

For protected routes, use middleware:

```go
func (h *handler) Register(router fiber.Router) {
    protected := router.Group("/", middleware.Auth)
    protected.Post("/", h.create)
    protected.Delete("/:id", h.delete)
}
```

## Handler Function Template

All handler functions that work with `c.UserContext()` should start a span and keep the traced context on Fiber:

```go
func (h *handler) getItems(c *fiber.Ctx) error {
    tracedCtx, span := tracing.StartSpan(c.UserContext(), "internal:framework:primary:http:item:get_items:getItems")
    defer span.End()
    c.SetUserContext(tracedCtx)

    // ...
}
```

Span names must use the file path plus function name format:

- `internal:framework:primary:http:<module>:<file_without_extension>:<function>`

### List (GET /)

```go
func (h *handler) getItems(c *fiber.Ctx) error {
    var (
        ctx = c.UserContext()
        req = new(restentity.GetItemsReq)
        v   = adapter.Adapters.Validator
    )

    if err := c.QueryParser(req); err != nil {
        log.Ctx(ctx).Warn().Err(err).Msg("Failed to parse query params")
        return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
    }

    req.SetDefault()

    if err := v.Validate(req); err != nil {
        log.Ctx(ctx).Warn().Err(err).Msg("Failed to validate query params")
        code, errors := errmsg.Errors(err, req)
        return c.Status(code).JSON(response.Error(errors))
    }

    filter := coreentity.ItemListFilter{
        TenantID: common.GetUserContext(ctx).TenantID,
        Q:        req.Q,
        Page:     req.Page,
        Paginate: req.Paginate,
    }

    items, total, err := h.core.GetItems(ctx, filter)
    if err != nil {
        log.Ctx(ctx).Error().Err(err).Msg("Failed to get items")
        code, errors := errmsg.Errors[error](err)
        return c.Status(code).JSON(response.Error(errors))
    }

    restItems := make([]restentity.Item, 0, len(items))
    for _, item := range items {
        restItems = append(restItems, mapper.ItemFromCoreToRest(item))
    }

    resp := &restentity.GetItemsResp{
        Items: restItems,
        Meta: types.Meta{
            Page:      req.Page,
            Paginate:  req.Paginate,
            TotalData: total,
        },
    }
    resp.Meta.CountTotalPage(req.Page, req.Paginate, total)

    return c.Status(fiber.StatusOK).JSON(response.Success(resp, ""))
}
```

### Get Single (GET /:id)

```go
func (h *handler) getItem(c *fiber.Ctx) error {
    var (
        ctx = c.UserContext()
        req = new(restentity.GetItemReq)
        v   = adapter.Adapters.Validator
    )

    if err := c.ParamsParser(req); err != nil { ... }
    if err := v.Validate(req); err != nil { ... }

    filter := coreentity.Item{
        TenantID: common.GetUserContext(ctx).TenantID,
        ID:       req.ID,
    }

    item, err := h.core.GetItem(ctx, filter)
    if err != nil { ... }

    return c.Status(fiber.StatusOK).JSON(response.Success(restentity.GetItemResp{
        Item: mapper.ItemFromCoreToRest(*item),
    }, ""))
}
```

### Create (POST /)

```go
func (h *handler) createItem(c *fiber.Ctx) error {
    var (
        ctx = c.UserContext()
        req = new(restentity.CreateItemReq)
        v   = adapter.Adapters.Validator
    )

    if err := c.BodyParser(req); err != nil { ... }
    if err := v.Validate(req); err != nil { ... }

    data := mapper.ItemFromRestCreateToCore(ctx, *req)
    created, err := h.core.CreateItem(ctx, data)
    if err != nil { ... }

    return c.Status(fiber.StatusCreated).JSON(response.Success(restentity.CreateItemResp{ID: created.ID}, ""))
}
```

### Update (PUT /:id)

```go
func (h *handler) updateItem(c *fiber.Ctx) error {
    var (
        ctx = c.UserContext()
        req = new(restentity.UpdateItemReq)
        v   = adapter.Adapters.Validator
    )

    if err := c.ParamsParser(req); err != nil { ... }
    if err := c.BodyParser(req); err != nil { ... }
    if err := v.Validate(req); err != nil { ... }

    data := mapper.ItemFromRestUpdateToCore(ctx, *req)
    if err := h.core.UpdateItem(ctx, data); err != nil { ... }

    return c.Status(fiber.StatusOK).JSON(response.Success(nil, ""))
}
```

### Delete (DELETE /:id)

```go
func (h *handler) deleteItem(c *fiber.Ctx) error {
    var (
        ctx = c.UserContext()
        req = new(restentity.DeleteItemReq)
        v   = adapter.Adapters.Validator
    )

    if err := c.ParamsParser(req); err != nil { ... }
    if err := v.Validate(req); err != nil { ... }

    filter := coreentity.ItemDeleteFilter{
        TenantID: common.GetUserContext(ctx).TenantID,
        ID:       req.ID,
    }

    if err := h.core.DeleteItem(ctx, filter); err != nil { ... }

    return c.Status(fiber.StatusOK).JSON(response.Success(nil, ""))
}
```

## Pagination

List endpoints use `types.MetaQuery` for pagination input and `types.Meta` for response:

### Request Struct

```go
type GetItemsReq struct {
    Q string `query:"q" validate:"omitempty,min=2"`
    types.MetaQuery  // embeds Page and Paginate
}

func (r *GetItemsReq) SetDefault() {
    r.MetaQuery.SetDefault()  // defaults: page=1, paginate=25
}
```

### Response Struct

```go
type GetItemsResp struct {
    Items []Item     `json:"items"`
    Meta  types.Meta `json:"meta"`
}
```

### Response JSON

```json
{
    "success": true,
    "message": "...",
    "data": {
        "items": [...],
        "meta": {
            "page": 1,
            "paginate": 25,
            "total_data": 100,
            "total_page": 4
        }
    }
}
```

## Log Level Rules in Handler

| Situation | Level | Example |
|-----------|-------|---------|
| Parse/validation error (client fault) | `Warn` | `log.Ctx(ctx).Warn().Err(err).Msg(...)` |
| Core/repo error (server fault) | `Error` | `log.Ctx(ctx).Error().Err(err).Msg(...)` |

package handler

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/entity/mapper"
	"codebase-app/internal/entity/restentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v3"
)

func (h *internalUserHandler) createInternalUser(c fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.Context(), "internal:framework:primary:http:internaluser:handler:createInternalUser")
	defer span.End()
	c.SetContext(tracedCtx)

	var (
		ctx = c.Context()
		req = new(restentity.CreateInternalUserReq)
		v   = adapter.Adapters.Validator
	)
	if err := c.Bind().Body(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}
	if err := v.Validate(req); err != nil {
		code, errs := errmsg.Errors(ctx, err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	item, err := h.core.CreateInternalUser(ctx, mapper.InternalUserFromRestCreateToCore(ctx, *req))
	if err != nil {
		code, errs := errmsg.Errors[error](ctx, err)
		return c.Status(code).JSON(response.Error(errs))
	}
	return c.Status(fiber.StatusCreated).JSON(response.Success(restentity.CreateInternalUserResp{ID: item.ID}, ""))
}

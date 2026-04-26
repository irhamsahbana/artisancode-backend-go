package handler

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/entity/mapper"
	"codebase-app/internal/entity/restentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v2"
)

func (h *internalQuotationHandler) createQuotation(c *fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.UserContext(), "internal:framework:primary:http:internalquotation:create_quotation:createQuotation")
	defer span.End()
	c.SetUserContext(tracedCtx)

	ctx := c.UserContext()
	req := new(restentity.CreateInternalQuotationReq)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}
	if err := adapter.Adapters.Validator.Validate(req); err != nil {
		code, errors := errmsg.Errors(ctx, err, req)
		return c.Status(code).JSON(response.Error(errors))
	}
	item, err := h.core.CreateQuotation(ctx, mapper.CreateInternalQuotationFromRest(ctx, *req))
	if err != nil {
		code, errors := errmsg.Errors[error](ctx, err)
		return c.Status(code).JSON(response.Error(errors))
	}
	return c.Status(fiber.StatusCreated).JSON(response.Success(mapper.InternalQuotationToRest(*item), ""))
}

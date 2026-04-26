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

func (h *internalInvoiceHandler) getPaymentAttempts(c *fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.UserContext(), "internal:framework:primary:http:internalinvoice:get_payment_attempts:getPaymentAttempts")
	defer span.End()
	c.SetUserContext(tracedCtx)

	ctx := c.UserContext()
	req := new(restentity.GetInternalCommerceResourceReq)
	if err := c.ParamsParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}
	if err := adapter.Adapters.Validator.Validate(req); err != nil {
		code, errors := errmsg.Errors(ctx, err, req)
		return c.Status(code).JSON(response.Error(errors))
	}
	items, err := h.core.GetPaymentAttempts(ctx, req.ID)
	if err != nil {
		code, errors := errmsg.Errors[error](ctx, err)
		return c.Status(code).JSON(response.Error(errors))
	}
	resp := restentity.GetPaymentAttemptsResp{Items: make([]restentity.InternalPaymentAttempt, 0, len(items))}
	for _, item := range items {
		resp.Items = append(resp.Items, mapper.InternalPaymentAttemptToRest(item, nil))
	}
	return c.Status(fiber.StatusOK).JSON(response.Success(resp, ""))
}

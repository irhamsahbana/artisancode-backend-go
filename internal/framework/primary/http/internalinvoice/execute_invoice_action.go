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

func (h *internalInvoiceHandler) executeInvoiceAction(c fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.Context(), "internal:framework:primary:http:internalinvoice:execute_invoice_action:executeInvoiceAction")
	defer span.End()
	c.SetContext(tracedCtx)

	ctx := c.Context()
	req := new(restentity.ExecuteInvoiceActionReq)
	if err := c.Bind().URI(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}
	if err := c.Bind().Body(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}
	if err := adapter.Adapters.Validator.Validate(req); err != nil {
		code, errors := errmsg.Errors(ctx, err, req)
		return c.Status(code).JSON(response.Error(errors))
	}
	result, err := h.core.ExecuteInvoiceAction(ctx, mapper.InternalInvoiceActionFromRest(ctx, *req, c.Get("X-Request-ID")))
	if err != nil {
		code, errors := errmsg.Errors[error](ctx, err)
		return c.Status(code).JSON(response.Error(errors))
	}
	return c.Status(fiber.StatusOK).JSON(response.Success(mapper.InternalPaymentAttemptToRest(result.PaymentAttempt, result.Instruction), ""))
}

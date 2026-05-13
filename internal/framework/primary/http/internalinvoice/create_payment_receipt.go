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

func (h *internalInvoiceHandler) createPaymentReceipt(c fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.Context(), "internal:framework:primary:http:internalinvoice:create_payment_receipt:createPaymentReceipt")
	defer span.End()
	c.SetContext(tracedCtx)

	ctx := c.Context()
	req := new(restentity.CreatePaymentReceiptReq)
	if err := c.Bind().Body(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}
	if err := adapter.Adapters.Validator.Validate(req); err != nil {
		code, errors := errmsg.Errors(ctx, err, req)
		return c.Status(code).JSON(response.Error(errors))
	}
	result, err := h.core.CreatePaymentReceipt(ctx, mapper.InternalPaymentReceiptFromRest(ctx, *req))
	if err != nil {
		code, errors := errmsg.Errors[error](ctx, err)
		return c.Status(code).JSON(response.Error(errors))
	}
	return c.Status(fiber.StatusCreated).JSON(response.Success(restentity.CreatePaymentReceiptResp{
		ID: result.Receipt.ID, ReceiptNumber: result.Receipt.ReceiptNumber, Status: result.Receipt.Status,
		Invoice: restentity.InternalInvoiceSummary{ID: result.Invoice.ID, InvoiceNumber: result.Invoice.InvoiceNumber, Status: result.Invoice.Status},
		Order:   mapper.InternalOrderToRest(result.Order, &result.Invoice),
	}, ""))
}

package handler

import (
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v3"
)

func (h *internalBillingHandler) createManualInvoice(c fiber.Ctx) error {
	ctx, span := tracing.StartSpan(c.Context(), "internal:framework:primary:http:internalbilling:createManualInvoice")
	defer span.End()

	var input coreentity.InternalBillingManualInvoiceInput
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(400).JSON(response.Error(err.Error()))
	}

	result, err := h.core.CreateManualInvoice(ctx, input)
	if err != nil {
		return err
	}

	return c.Status(201).JSON(response.Success(result, ""))
}

func (h *internalBillingHandler) createPaymentReceipt(c fiber.Ctx) error {
	ctx, span := tracing.StartSpan(c.Context(), "internal:framework:primary:http:internalbilling:createPaymentReceipt")
	defer span.End()

	var input coreentity.InternalPaymentReceipt
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(400).JSON(response.Error(err.Error()))
	}

	result, err := h.core.CreatePaymentReceipt(ctx, input)
	if err != nil {
		return err
	}

	return c.Status(201).JSON(response.Success(result, ""))
}

func (h *internalBillingHandler) executePaymentReceiptAction(c fiber.Ctx) error {
	ctx, span := tracing.StartSpan(c.Context(), "internal:framework:primary:http:internalbilling:executePaymentReceiptAction")
	defer span.End()

	var input coreentity.InternalBillingPaymentReceiptActionInput
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(400).JSON(response.Error(err.Error()))
	}
	input.ID = c.Params("id")

	result, err := h.core.ExecutePaymentReceiptAction(ctx, input)
	if err != nil {
		return err
	}

	return c.JSON(response.Success(result, ""))
}

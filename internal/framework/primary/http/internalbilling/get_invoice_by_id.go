package handler

import (
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v3"
)

func (h *internalBillingHandler) getInvoiceByID(c fiber.Ctx) error {
	ctx, span := tracing.StartSpan(c.Context(), "internal:framework:primary:http:internalbilling:getInvoiceByID")
	defer span.End()

	id := c.Params("id")

	result, err := h.core.GetInvoiceByID(ctx, id)
	if err != nil {
		return err
	}

	return c.JSON(response.Success(result, ""))
}

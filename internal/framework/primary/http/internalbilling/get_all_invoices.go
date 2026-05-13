package handler

import (
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v3"
)

func (h *internalBillingHandler) getAllInvoices(c fiber.Ctx) error {
	ctx, span := tracing.StartSpan(c.Context(), "internal:framework:primary:http:internalbilling:getAllInvoices")
	defer span.End()

	var filter coreentity.InternalBillingInvoiceListFilter
	if v := c.Query("tenant_id"); v != "" {
		filter.TenantID = &v
	}
	filter.Status = c.Query("status")
	filter.FromDate = c.Query("from_date")
	filter.ToDate = c.Query("to_date")

	result, err := h.core.GetAllInvoices(ctx, filter)
	if err != nil {
		return err
	}

	return c.JSON(response.Success(result, ""))
}

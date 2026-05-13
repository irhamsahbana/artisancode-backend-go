package handler

import (
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v3"
)

func (h *internalBillingHandler) getReconciliationCases(c fiber.Ctx) error {
	ctx, span := tracing.StartSpan(c.Context(), "internal:framework:primary:http:internalbilling:getReconciliationCases")
	defer span.End()

	filter := coreentity.InternalBillingReconciliationCaseFilter{
		Status: c.Query("status"),
	}

	result, err := h.core.GetReconciliationCases(ctx, filter)
	if err != nil {
		return err
	}

	return c.JSON(response.Success(result, ""))
}

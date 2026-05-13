package handler

import (
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v3"
)

func (h *internalBillingHandler) getLedgerEntries(c fiber.Ctx) error {
	ctx, span := tracing.StartSpan(c.Context(), "internal:framework:primary:http:internalbilling:getLedgerEntries")
	defer span.End()

	tenantID := c.Params("tenantId")
	filter := coreentity.InternalBillingLedgerListFilter{
		EntryTypes: parseStringSlice(c.Query("entry_types")),
		FromDate:   c.Query("from_date"),
		ToDate:     c.Query("to_date"),
	}

	result, err := h.core.GetLedgerEntries(ctx, tenantID, filter)
	if err != nil {
		return err
	}

	return c.JSON(response.Success(result, ""))
}

func parseStringSlice(s string) []string {
	if s == "" {
		return nil
	}
	var out []string
	start := 0
	for i, c := range s {
		if c == ',' {
			out = append(out, s[start:i])
			start = i + 1
		}
	}
	out = append(out, s[start:])
	return out
}

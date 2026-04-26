package handler

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/restentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"
	"codebase-app/pkg/types"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *internalInvoiceHandler) getInvoices(c *fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.UserContext(), "internal:framework:primary:http:internalinvoice:get_invoices:getInvoices")
	defer span.End()
	c.SetUserContext(tracedCtx)

	var (
		ctx = c.UserContext()
		req = new(restentity.GetInternalCommerceResourcesReq)
		v   = adapter.Adapters.Validator
	)

	if err := c.QueryParser(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to parse internal invoices query params")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}
	req.SetDefault()

	if err := v.Validate(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to validate internal invoices query params")
		code, errors := errmsg.Errors(ctx, err, req)
		return c.Status(code).JSON(response.Error(errors))
	}

	items, total, err := h.core.GetInvoices(ctx, coreentity.InternalCommerceListFilter{
		Page:     req.Page,
		Paginate: req.Paginate,
	})
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get internal invoices")
		code, errors := errmsg.Errors[error](ctx, err)
		return c.Status(code).JSON(response.Error(errors))
	}

	restItems := make([]restentity.InternalInvoice, 0, len(items))
	for _, item := range items {
		restItems = append(restItems, restentity.InternalInvoice{
			ID: item.ID, InternalOrderID: item.InternalOrderID, InvoiceNumber: item.InvoiceNumber, Status: item.Status,
			CurrencyCode: item.CurrencyCode, Amount: item.Amount.StringFixed(2), AmountPaid: item.AmountPaid.StringFixed(2),
			AmountOutstanding: item.AmountOutstanding.StringFixed(2), DueAt: item.DueAt, PaidAt: item.PaidAt,
			ExpiredAt: item.ExpiredAt, Metadata: item.Metadata, CreatedAt: item.CreatedAt, UpdatedAt: item.UpdatedAt,
		})
	}

	resp := &restentity.GetInternalInvoicesResp{
		Items: restItems,
		Meta: types.Meta{
			Page:      req.Page,
			Paginate:  req.Paginate,
			TotalData: total,
		},
	}
	resp.Meta.CountTotalPage(req.Page, req.Paginate, total)

	return c.Status(fiber.StatusOK).JSON(response.Success(resp, ""))
}

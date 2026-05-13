package handler

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/mapper"
	"codebase-app/internal/entity/restentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

func (h *internalTenantBillingHandler) getInvoices(c fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(
		c.Context(),
		"internal:framework:primary:http:internaltenantbilling:get_invoices:getInvoices",
	)
	defer span.End()
	c.SetContext(tracedCtx)

	ctx := c.Context()
	req := new(restentity.GetTenantBillingInvoicesReq)
	if err := c.Bind().Query(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to parse tenant billing invoices query params")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}
	req.SetDefault()
	if err := adapter.Adapters.Validator.Validate(req); err != nil {
		code, errors := errmsg.Errors(ctx, err, req)
		return c.Status(code).JSON(response.Error(errors))
	}

	items, err := h.core.GetInvoices(ctx, coreentity.TenantBillingInvoiceListFilter{
		UserCtx:  common.GetUserContext(ctx),
		Page:     req.Page,
		Paginate: req.Paginate,
	})
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get tenant billing invoices")
		code, errors := errmsg.Errors[error](ctx, err)
		return c.Status(code).JSON(response.Error(errors))
	}

	restItems := make([]restentity.TenantBillingInvoice, 0, len(items))
	for _, item := range items {
		restItems = append(restItems, mapper.TenantBillingInvoiceFromCoreToRest(item))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(restentity.GetTenantBillingInvoicesResp{
		Items: restItems,
	}, ""))
}

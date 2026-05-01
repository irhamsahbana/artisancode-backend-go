package core

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

func (c *internalQuotationCore) executeApproveQuotationAction(
	ctx context.Context,
	input coreentity.InternalQuotationActionInput,
) (*coreentity.InternalCommerceBundle, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:core:internalquotation:approve_quotation_action:executeApproveQuotationAction",
	)
	defer span.End()

	quote, err := c.quotationRepo.GetQuotation(ctx, input.UserCtx.TenantID, input.ID)
	if err != nil {
		return nil, err
	}
	if quote.Status != coreentity.QuotationStatusDraft && quote.Status != coreentity.QuotationStatusSent {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, input).Msg(errmsg.MessageQuotationCannotBeApprovedFromCurrentStatus)
		return nil, errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageQuotationCannotBeApprovedFromCurrentStatus)
	}
	approved, err := c.quotationRepo.ApproveQuotation(ctx, input.UserCtx.TenantID, input.ID)
	if err != nil {
		return nil, err
	}
	return &coreentity.InternalCommerceBundle{Quotation: approved}, nil
}

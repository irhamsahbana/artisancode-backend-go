package core

import (
	"context"
	"strings"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

func (c *internalQuotationCore) ExecuteQuotationAction(
	ctx context.Context,
	input coreentity.InternalQuotationActionInput,
) (*coreentity.InternalCommerceBundle, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:core:internalquotation:execute_quotation_action:ExecuteQuotationAction",
	)
	defer span.End()

	input.Action = strings.TrimSpace(input.Action)
	switch input.Action {
	case coreentity.ActionApproveQuotation:
		return c.executeApproveQuotationAction(ctx, input)
	case coreentity.ActionConvertQuotation:
		return c.executeConvertQuotationAction(ctx, input)
	default:
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, input).Msg("Unsupported quotation action")
		return nil, errmsg.NewCustomErrors(400).SetMessage("Unsupported quotation action")
	}
}

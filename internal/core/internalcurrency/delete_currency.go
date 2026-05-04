package core

import (
	"context"
	"strings"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *internalCurrencyCore) DeleteInternalCurrency(
	ctx context.Context,
	filter coreentity.InternalCurrencyDeleteFilter,
) error {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalcurrency:delete_currency:DeleteInternalCurrency")
	defer span.End()

	if err := c.authorizeWrite(common.GetUserContext(ctx)); err != nil {
		return err
	}

	filter.Code = strings.ToUpper(strings.TrimSpace(filter.Code))
	return c.repo.DeleteInternalCurrency(ctx, filter)
}

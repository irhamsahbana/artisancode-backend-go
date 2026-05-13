package core

import (
	"context"
	"fmt"

	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
)

func (c *internalTenantBillingCore) CheckUsageLimit(
	ctx context.Context,
	resourceType string,
	currentCount int64,
) error {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internaltenantbilling:check_usage_limit:CheckUsageLimit")
	defer span.End()

	snapshot, err := c.GetEntitlements(ctx)
	if err != nil {
		return err
	}

	limit, ok := snapshot.UsageLimits[resourceType]
	if !ok {
		return nil
	}

	if limit == 0 {
		return nil
	}

	if currentCount >= limit {
		return errmsg.NewCustomErrors(422,
			errmsg.WithMessage(fmt.Sprintf(
				"usage limit exceeded for %s: %d/%d",
				resourceType, currentCount, limit,
			)),
		).SetErrorCode(errmsg.MessageUsageLimitExceeded)
	}

	return nil
}

package core

import (
	"context"
	"time"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *messageQueueCore) CleanupProcessedMessages(ctx context.Context, req coreentity.CleanupProcessedMessageQueueReq) (*coreentity.CleanupProcessedMessageQueueResp, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:message_queue:cleanup_processed_messages:CleanupProcessedMessages")
	defer span.End()

	before := req.Before
	if before == "" {
		before = time.Now().UTC().Format(time.RFC3339)
	}
	if req.Limit <= 0 {
		req.Limit = 100
	}

	return c.repo.CleanupProcessedMessages(ctx, coreentity.CleanupProcessedMessageQueueReq{
		Before: before,
		Limit:  req.Limit,
	})
}

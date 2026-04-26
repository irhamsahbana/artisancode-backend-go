package repository

import (
	"context"
	"time"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *messageQueueRepo) CleanupProcessedMessages(
	ctx context.Context,
	req coreentity.CleanupProcessedMessageQueueReq,
) (*coreentity.CleanupProcessedMessageQueueResp, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:message_queue:cleanup_processed_messages:CleanupProcessedMessages",
	)
	defer span.End()

	before, err := time.Parse(time.RFC3339, req.Before)
	if err != nil {
		log.Ctx(ctx).
			Error().
			Err(err).
			Any(common.LogKeyPayload, req).
			Msg("Failed to parse processed message cleanup time")
		return nil, err
	}

	if req.Limit <= 0 {
		req.Limit = 100
	}

	query := `
		WITH deleted_rows AS (
			DELETE FROM message_queue
			WHERE id IN (
				SELECT id
				FROM message_queue
				WHERE deleted_at IS NULL
					AND status = 'processed'
					AND processed_at IS NOT NULL
					AND processed_at < ?
				ORDER BY processed_at ASC
				LIMIT ?
			)
			RETURNING id
		)
		SELECT COUNT(*) FROM deleted_rows
	`

	var deleted int
	err = r.db.GetContext(ctx, &deleted, r.db.Rebind(query), before, req.Limit)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, req).Msg("Failed to clean processed message queue rows")
		return nil, err
	}

	return &coreentity.CleanupProcessedMessageQueueResp{
		Deleted: deleted,
	}, nil
}

package repository

import (
	"context"
	"fmt"
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

	topics := []struct {
		messageTable string
		offsetTable  string
	}{
		{`"watermill_email_verification"`, `"watermill_offsets_email_verification"`},
		{`"watermill_email_forgot_password"`, `"watermill_offsets_email_forgot_password"`},
		{`"watermill_email_invitation"`, `"watermill_offsets_email_invitation"`},
		{`"watermill_export_job_requested"`, `"watermill_offsets_export_job_requested"`},
		{`"watermill_dead_letter"`, `"watermill_offsets_dead_letter"`},
	}

	var deleted int
	for _, topic := range topics {
		remainingLimit := req.Limit - deleted
		if remainingLimit <= 0 {
			break
		}

		query := fmt.Sprintf(`
			WITH cutoff AS (
				SELECT
					MIN(last_processed_transaction_id) AS transaction_id,
					MIN(offset_acked) AS offset_acked
				FROM %s
			),
			deleted_rows AS (
				DELETE FROM %s messages
				WHERE messages.created_at < ?
					AND EXISTS (SELECT 1 FROM cutoff WHERE transaction_id IS NOT NULL)
					AND (
						messages.transaction_id < (SELECT transaction_id FROM cutoff)
						OR (
							messages.transaction_id = (SELECT transaction_id FROM cutoff)
							AND messages."offset" <= (SELECT offset_acked FROM cutoff)
						)
					)
					AND messages."offset" IN (
						SELECT selected."offset"
						FROM %s selected
						WHERE selected.created_at < ?
						ORDER BY selected.transaction_id ASC, selected."offset" ASC
						LIMIT ?
					)
				RETURNING messages."offset"
			)
			SELECT COUNT(*) FROM deleted_rows
		`, topic.offsetTable, topic.messageTable, topic.messageTable)

		var topicDeleted int
		err = r.db.GetContext(ctx, &topicDeleted, r.db.Rebind(query), before, before, remainingLimit)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, req).Msg("Failed to clean processed Watermill message rows")
			return nil, err
		}

		deleted += topicDeleted
	}

	return &coreentity.CleanupProcessedMessageQueueResp{
		Deleted: deleted,
	}, nil
}

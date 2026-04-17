package repository

import (
	"context"

	"codebase-app/internal/entity/coreentity"
)

type MessageQueueRepository interface {
	CleanupProcessedMessages(ctx context.Context, req coreentity.CleanupProcessedMessageQueueReq) (*coreentity.CleanupProcessedMessageQueueResp, error)
}

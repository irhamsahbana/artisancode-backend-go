package core

import (
	"codebase-app/internal/entity/coreentity"
	"context"
)

type WebhookCore interface {
	HandleDOKUWebhook(ctx context.Context, notification coreentity.DOKUWebhookNotification) (*coreentity.DOKUWebhookResult, error)
}

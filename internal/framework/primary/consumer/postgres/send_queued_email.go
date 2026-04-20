package consumer

import (
	"context"
	"fmt"

	emailint "codebase-app/internal/integration/email"
)

func sendQueuedEmail(ctx context.Context, payload emailint.EmailPayload) error {
	sender := emailint.GetEmailSender()
	if sender == nil {
		return fmt.Errorf("email sender not initialized")
	}

	return sender.SendSync(payload)
}

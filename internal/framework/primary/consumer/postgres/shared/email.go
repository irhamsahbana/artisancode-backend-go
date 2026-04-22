package shared

import (
	"context"
	"fmt"

	infraTracing "codebase-app/internal/infrastructure/tracing"
	emailint "codebase-app/internal/integration/email"
)

func SendQueuedEmail(ctx context.Context, payload emailint.EmailPayload) error {
	ctx, span := infraTracing.StartSpan(ctx, "internal:framework:primary:consumer:postgres:shared:email:SendQueuedEmail")
	defer span.End()

	if err := ctx.Err(); err != nil {
		infraTracing.RecordError(span, err)
		return err
	}

	sender := emailint.GetEmailSender()
	if sender == nil {
		err := fmt.Errorf("email sender not initialized")
		infraTracing.RecordError(span, err)
		return err
	}

	err := sender.SendSync(ctx, payload)
	if err != nil {
		infraTracing.RecordError(span, err)
	}

	return err
}

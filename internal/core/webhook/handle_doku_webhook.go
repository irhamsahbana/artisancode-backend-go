package core

import (
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"context"
	"errors"

	"github.com/rs/zerolog/log"
)

func (c *webhookCore) HandleDOKUWebhook(
	ctx context.Context,
	notification coreentity.DOKUWebhookNotification,
) (*coreentity.DOKUWebhookResult, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:webhook:handle_doku_webhook:HandleDOKUWebhook")
	defer span.End()

	if c.dokuVerifier == nil {
		err := errors.New("doku webhook verifier is required")
		tracing.RecordError(span, err)
		return nil, err
	}

	if !c.dokuVerifier.VerifyWebhookSignatureHeaders(
		notification.Headers,
		notification.RawBody,
		notification.TargetPath,
	) {
		err := errors.New("invalid doku webhook signature")
		log.Ctx(ctx).Warn().
			Str("invoice_number", notification.Event.Order.InvoiceNumber).
			Str("request_id", notification.Headers.RequestID).
			Msg("DOKU webhook signature verification failed")
		tracing.RecordError(span, err)
		return nil, err
	}

	log.Ctx(ctx).Info().
		Str("invoice_number", notification.Event.Order.InvoiceNumber).
		Str("order_status", notification.Event.Order.Status).
		Str("payment_status", notification.Event.Transaction.Status).
		Str("request_id", notification.Headers.RequestID).
		Msg("DOKU webhook accepted")

	return &coreentity.DOKUWebhookResult{
		Provider:      "doku",
		InvoiceNumber: notification.Event.Order.InvoiceNumber,
		OrderStatus:   notification.Event.Order.Status,
		PaymentStatus: notification.Event.Transaction.Status,
	}, nil
}

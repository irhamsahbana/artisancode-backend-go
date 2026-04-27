package user

import (
	"context"
	"encoding/json"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/framework/primary/consumer/postgres/shared"
	"codebase-app/internal/infrastructure/tracing"
	emailint "codebase-app/internal/integration/email"

	watermillMessage "github.com/ThreeDotsLabs/watermill/message"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
)

func (a *Adapter) ForgotPasswordHandler(msg *watermillMessage.Message) error {
	ctx := otel.GetTextMapPropagator().Extract(
		msg.Context(),
		shared.WatermillMetadataCarrier(msg.Metadata),
	)
	ctx, span := tracing.StartSpan(ctx, "internal:framework:primary:consumer:postgres:user:forgot_password:ForgotPasswordHandler")
	defer span.End()

	payload := &coreentity.QueuedEmailMessage{}

	err := json.Unmarshal(msg.Payload, payload)
	if err != nil {
		tracing.RecordError(span, err)
		log.Ctx(ctx).Error().Err(err).Msg("failed to unmarshal forgot password message")
		msg.Ack()
		return nil
	}

	logger := shared.LogEmailMessage(ctx, msg.Metadata.Get("topic"), 1)
	logger.Debug().Msg("received forgot password message")

	err = a.sendForgotPassword(ctx, payload)
	if err != nil {
		tracing.RecordError(span, err)
		logger.Error().Err(err).Msg("failed to send forgot password email")
		return err
	}

	logger.Info().Msg("finished forgot password message")
	return nil
}

func (a *Adapter) sendForgotPassword(ctx context.Context, payload *coreentity.QueuedEmailMessage) error {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:primary:consumer:postgres:user:forgot_password:sendForgotPassword")
	defer span.End()

	rendered, err := emailint.BuildPasswordResetEmail(emailint.AuthTemplateInput{
		UserName:          shared.CoalesceEmailIdentity(payload.UserName, payload.Email),
		TenantName:        payload.TenantName,
		ActionLink:        payload.ActionLink,
		PreferredLanguage: payload.PreferredLanguage,
	})
	if err != nil {
		tracing.RecordError(span, err)
		log.Ctx(ctx).Error().Err(err).Msg("failed to render password reset template")
		return err
	}

	return shared.SendQueuedEmail(ctx, emailint.EmailPayload{
		To:      []string{payload.Email},
		Subject: rendered.Subject,
		Body:    rendered.Body,
		IsHTML:  true,
	})
}

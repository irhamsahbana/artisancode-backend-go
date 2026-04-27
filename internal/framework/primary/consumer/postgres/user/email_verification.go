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

func (a *Adapter) EmailVerificationHandler(msg *watermillMessage.Message) error {
	ctx := otel.GetTextMapPropagator().Extract(
		msg.Context(),
		shared.WatermillMetadataCarrier(msg.Metadata),
	)
	ctx, span := tracing.StartSpan(ctx, "internal:framework:primary:consumer:postgres:user:email_verification:EmailVerificationHandler")
	defer span.End()

	payload := &coreentity.QueuedEmailMessage{}

	err := json.Unmarshal(msg.Payload, payload)
	if err != nil {
		tracing.RecordError(span, err)
		log.Ctx(ctx).Error().Err(err).Msg("failed to unmarshal email verification message")
		msg.Ack()
		return nil
	}

	logger := shared.LogEmailMessage(ctx, msg.Metadata.Get("topic"), 1)
	logger.Debug().Msg("received verification message")

	err = a.sendEmailVerification(ctx, payload)
	if err != nil {
		tracing.RecordError(span, err)
		logger.Error().Err(err).Msg("failed to send verification email")
		return err
	}

	logger.Info().Msg("finished verification message")
	return nil
}

func (a *Adapter) sendEmailVerification(ctx context.Context, payload *coreentity.QueuedEmailMessage) error {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:primary:consumer:postgres:user:email_verification:sendEmailVerification")
	defer span.End()

	rendered, err := emailint.BuildVerificationEmail(emailint.AuthTemplateInput{
		UserName:          shared.CoalesceEmailIdentity(payload.UserName, payload.Email),
		TenantName:        payload.TenantName,
		ActionLink:        payload.ActionLink,
		PreferredLanguage: payload.PreferredLanguage,
	})
	if err != nil {
		tracing.RecordError(span, err)
		log.Ctx(ctx).Error().Err(err).Msg("failed to render verification template")
		return err
	}

	return shared.SendQueuedEmail(ctx, emailint.EmailPayload{
		To:      []string{payload.Email},
		Subject: rendered.Subject,
		Body:    rendered.Body,
		IsHTML:  true,
	})
}

package user

import (
	"context"
	"encoding/json"

	"codebase-app/internal/framework/primary/consumer/postgres/shared"
	"codebase-app/internal/infrastructure/tracing"
	emailint "codebase-app/internal/integration/email"
	integrationPorts "codebase-app/internal/ports/secondary/integration"

	"github.com/rs/zerolog/log"
)

type EmailVerificationEventPayload struct {
	Email             string `json:"email"`
	UserName          string `json:"user_name"`
	TenantName        string `json:"tenant_name"`
	ActionLink        string `json:"action_link"`
	PreferredLanguage string `json:"preferred_language"`
}

func EmailVerificationHandler(ctx context.Context, msg integrationPorts.MessageBusMessage) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:primary:consumer:postgres:user:email_verification:EmailVerificationHandler")
	defer span.End()

	payload := &EmailVerificationEventPayload{}

	err := json.Unmarshal(msg.Data(), payload)
	if err != nil {
		tracing.RecordError(span, err)
		log.Ctx(ctx).Error().Err(err).Msg("failed to unmarshal email verification message")
		if err := msg.Ack(ctx); err != nil {
			tracing.RecordError(span, err)
			log.Ctx(ctx).Error().Err(err).Msg("failed to acknowledge invalid verification message")
		}
		return
	}

	logger := shared.LogEmailMessage(ctx, msg.Subject(), 1)
	logger.Debug().Msg("received verification message")

	err = sendEmailVerification(ctx, payload)
	if err != nil {
		sendErr := err
		tracing.RecordError(span, err)
		logger.Error().Err(err).Msg("failed to send verification email")
		if nakErr := msg.Nak(ctx, sendErr.Error()); nakErr != nil {
			tracing.RecordError(span, nakErr)
			logger.Error().Err(nakErr).Msg("failed to negative-acknowledge verification message")
		}
		return
	}
	if err := msg.Ack(ctx); err != nil {
		tracing.RecordError(span, err)
		logger.Error().Err(err).Msg("failed to acknowledge verification message")
		return
	}

	logger.Info().Msg("finished verification message")
}

func sendEmailVerification(ctx context.Context, payload *EmailVerificationEventPayload) error {
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

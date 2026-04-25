package user

import (
	"context"
	"encoding/json"

	"codebase-app/internal/framework/primary/consumer/postgres/shared"
	"codebase-app/internal/infrastructure/tracing"
	emailint "codebase-app/internal/integration/email"
	integrationPorts "codebase-app/internal/ports/integration"

	"github.com/rs/zerolog/log"
)

type ForgotPasswordEventPayload struct {
	Email             string `json:"email"`
	UserName          string `json:"user_name"`
	TenantName        string `json:"tenant_name"`
	ActionLink        string `json:"action_link"`
	PreferredLanguage string `json:"preferred_language"`
}

func ForgotPasswordHandler(ctx context.Context, msg integrationPorts.MessageBusMessage) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:primary:consumer:postgres:user:forgot_password:ForgotPasswordHandler")
	defer span.End()

	payload := &ForgotPasswordEventPayload{}

	err := json.Unmarshal(msg.Data(), payload)
	if err != nil {
		tracing.RecordError(span, err)
		log.Ctx(ctx).Error().Err(err).Msg("failed to unmarshal forgot password message")
		if err := msg.Ack(ctx); err != nil {
			tracing.RecordError(span, err)
			log.Ctx(ctx).Error().Err(err).Msg("failed to acknowledge invalid forgot password message")
		}
		return
	}

	logger := shared.LogEmailMessage(ctx, msg.Subject(), 1)
	logger.Debug().Msg("received forgot password message")

	err = sendForgotPassword(ctx, payload)
	if err != nil {
		sendErr := err
		tracing.RecordError(span, err)
		logger.Error().Err(err).Msg("failed to send forgot password email")
		if nakErr := msg.Nak(ctx, sendErr.Error()); nakErr != nil {
			tracing.RecordError(span, nakErr)
			logger.Error().Err(nakErr).Msg("failed to negative-acknowledge forgot password message")
		}
		return
	}

	if err := msg.Ack(ctx); err != nil {
		tracing.RecordError(span, err)
		logger.Error().Err(err).Msg("failed to acknowledge forgot password message")
		return
	}

	logger.Info().Msg("finished forgot password message")
}

func sendForgotPassword(ctx context.Context, payload *ForgotPasswordEventPayload) error {
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

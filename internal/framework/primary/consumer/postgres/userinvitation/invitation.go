package userinvitation

import (
	"context"
	"encoding/json"

	"codebase-app/internal/framework/primary/consumer/postgres/shared"
	"codebase-app/internal/infrastructure/tracing"
	emailint "codebase-app/internal/integration/email"
	integrationPorts "codebase-app/internal/ports/secondary/integration"

	"github.com/rs/zerolog/log"
)

type InvitationEventPayload struct {
	Email             string `json:"email"`
	UserName          string `json:"user_name"`
	TenantName        string `json:"tenant_name"`
	ActionLink        string `json:"action_link"`
	PreferredLanguage string `json:"preferred_language"`
}

func InvitationHandler(ctx context.Context, msg integrationPorts.MessageBusMessage) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:primary:consumer:postgres:userinvitation:invitation:InvitationHandler")
	defer span.End()

	payload := &InvitationEventPayload{}

	err := json.Unmarshal(msg.Data(), payload)
	if err != nil {
		tracing.RecordError(span, err)
		log.Ctx(ctx).Error().Err(err).Msg("failed to unmarshal invitation message")
		if err := msg.Ack(ctx); err != nil {
			tracing.RecordError(span, err)
			log.Ctx(ctx).Error().Err(err).Msg("failed to acknowledge invalid invitation message")
		}
		return
	}

	logger := shared.LogEmailMessage(ctx, msg.Subject(), 1)
	logger.Debug().Msg("received invitation message")

	emailPayload, err := buildInvitationEmailPayload(payload)
	if err == nil {
		err = shared.SendQueuedEmail(ctx, emailPayload)
	}
	if err != nil {
		sendErr := err
		tracing.RecordError(span, err)
		logger.Error().Err(err).Msg("failed to send invitation email")
		if nakErr := msg.Nak(ctx, sendErr.Error()); nakErr != nil {
			tracing.RecordError(span, nakErr)
			logger.Error().Err(nakErr).Msg("failed to negative-acknowledge invitation message")
		}
		return
	}

	if err := msg.Ack(ctx); err != nil {
		tracing.RecordError(span, err)
		logger.Error().Err(err).Msg("failed to acknowledge invitation message")
		return
	}

	logger.Info().Msg("finished invitation message")
}

func buildInvitationEmailPayload(payload *InvitationEventPayload) (emailint.EmailPayload, error) {
	rendered, err := emailint.BuildInvitationEmail(emailint.AuthTemplateInput{
		UserName:          shared.CoalesceEmailIdentity(payload.UserName, payload.Email),
		TenantName:        payload.TenantName,
		ActionLink:        payload.ActionLink,
		PreferredLanguage: payload.PreferredLanguage,
	})
	if err != nil {
		return emailint.EmailPayload{}, err
	}

	return emailint.EmailPayload{
		To:      []string{payload.Email},
		Subject: rendered.Subject,
		Body:    rendered.Body,
		IsHTML:  true,
	}, nil
}

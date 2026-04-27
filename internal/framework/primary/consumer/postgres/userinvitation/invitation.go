package userinvitation

import (
	"encoding/json"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/framework/primary/consumer/postgres/shared"
	"codebase-app/internal/infrastructure/tracing"
	emailint "codebase-app/internal/integration/email"

	watermillMessage "github.com/ThreeDotsLabs/watermill/message"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
)

func (a *Adapter) InvitationHandler(msg *watermillMessage.Message) error {
	ctx := otel.GetTextMapPropagator().Extract(
		msg.Context(),
		shared.WatermillMetadataCarrier(msg.Metadata),
	)
	ctx, span := tracing.StartSpan(ctx, "internal:framework:primary:consumer:postgres:userinvitation:invitation:InvitationHandler")
	defer span.End()

	payload := &coreentity.QueuedEmailMessage{}

	err := json.Unmarshal(msg.Payload, payload)
	if err != nil {
		tracing.RecordError(span, err)
		log.Ctx(ctx).Error().Err(err).Msg("failed to unmarshal invitation message")
		msg.Ack()
		return nil
	}

	logger := shared.LogEmailMessage(ctx, msg.Metadata.Get("topic"), 1)
	logger.Debug().Msg("received invitation message")

	emailPayload, err := a.buildInvitationEmailPayload(payload)
	if err == nil {
		err = shared.SendQueuedEmail(ctx, emailPayload)
	}
	if err != nil {
		tracing.RecordError(span, err)
		logger.Error().Err(err).Msg("failed to send invitation email")
		return err
	}

	logger.Info().Msg("finished invitation message")
	return nil
}

func (a *Adapter) buildInvitationEmailPayload(payload *coreentity.QueuedEmailMessage) (emailint.EmailPayload, error) {
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

package consumer

import (
	"codebase-app/internal/entity/common"
	infraTracing "codebase-app/internal/infrastructure/tracing"
	integrationPorts "codebase-app/internal/ports/secondary/integration"
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
)

func (a *App) Run(ctx context.Context) error {
	ctx, span := infraTracing.StartSpan(ctx, "consumer.Run")
	defer span.End()

	err := a.validate()
	if err != nil {
		infraTracing.RecordError(span, err)
		return err
	}

	defer func() {
		if err := a.exportPublisher.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close export message publisher")
		}
		if err := a.subscriptionManager.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close message subscription manager")
		}
	}()

	emailConsumeCtx, err := a.emailSubscription.Consume(func(msg integrationPorts.MessageBusMessage) {
		msgCtx := otel.GetTextMapPropagator().Extract(ctx, MessageHeadersCarrier(msg.Headers()))

		switch msg.Subject() {
		case common.MessageSubjectEmailVerification:
			EmailVerificationHandler(msgCtx, msg)
		case common.MessageSubjectEmailForgotPassword:
			ForgotPasswordHandler(msgCtx, msg)
		default:
		}
	})
	if err != nil {
		infraTracing.RecordError(span, err)
		return err
	}
	defer emailConsumeCtx.Stop()

	exportConsumeCtx, err := a.exportSubscription.Consume(ExportJobRequestedHandler(ctx, a.exportCore))
	if err != nil {
		emailConsumeCtx.Stop()
		infraTracing.RecordError(span, err)
		return err
	}
	defer exportConsumeCtx.Stop()

	err = a.waitForShutdown(ctx)
	if err != nil {
		infraTracing.RecordError(span, err)
	}

	return err
}

func (a *App) validate() error {
	if a.emailSubscription == nil {
		return errors.New("email subscription is required")
	}
	if a.exportSubscription == nil {
		return errors.New("export subscription is required")
	}
	if a.exportCore == nil {
		return errors.New("export job processor is required")
	}
	if a.exportPublisher == nil {
		return errors.New("export publisher is required")
	}
	if a.subscriptionManager == nil {
		return errors.New("message subscription manager is required")
	}

	return nil
}

package consumer

import (
	"codebase-app/internal/entity/common"
	exportjobconsumer "codebase-app/internal/framework/primary/consumer/postgres/export_job"
	sharedconsumer "codebase-app/internal/framework/primary/consumer/postgres/shared"
	userconsumer "codebase-app/internal/framework/primary/consumer/postgres/user"
	userinvitationconsumer "codebase-app/internal/framework/primary/consumer/postgres/userinvitation"
	infraTracing "codebase-app/internal/infrastructure/tracing"
	integrationPorts "codebase-app/internal/ports/integration"
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
)

func (a *App) Run(ctx context.Context) error {
	ctx, stopSignals := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stopSignals()

	ctx, span := infraTracing.StartSpan(ctx, "internal:framework:primary:consumer:postgres:run:Run")
	defer span.End()

	err := a.build(ctx)
	if err != nil {
		infraTracing.RecordError(span, err)
		return err
	}

	err = a.validate()
	if err != nil {
		infraTracing.RecordError(span, err)
		return err
	}

	emailConsumeCtx, err := a.emailSubscription.Consume(ctx, func(consumeCtx context.Context, msg integrationPorts.MessageBusMessage) {
		msgCtx := otel.GetTextMapPropagator().Extract(consumeCtx, sharedconsumer.MessageHeadersCarrier(msg.Headers()))

		switch msg.Subject() {
		case common.MessageSubjectEmailVerification:
			userconsumer.EmailVerificationHandler(msgCtx, msg)
		case common.MessageSubjectEmailForgotPassword:
			userconsumer.ForgotPasswordHandler(msgCtx, msg)
		case common.MessageSubjectEmailInvitation:
			userinvitationconsumer.InvitationHandler(msgCtx, msg)
		default:
		}
	})
	if err != nil {
		infraTracing.RecordError(span, err)
		return err
	}

	exportConsumeCtx, err := a.exportSubscription.Consume(ctx, exportjobconsumer.RequestedHandler(a.exportCore))
	if err != nil {
		emailConsumeCtx.Stop()
		infraTracing.RecordError(span, err)
		return err
	}

	err = a.waitForShutdown(ctx)
	if err != nil {
		infraTracing.RecordError(span, err)
	}

	exportConsumeCtx.Stop()
	emailConsumeCtx.Stop()

	if shutdownErr := a.shutdownResources(ctx); shutdownErr != nil {
		if err == nil {
			err = shutdownErr
		} else {
			log.Ctx(ctx).Error().Err(shutdownErr).Msg("failed to shutdown resources cleanly")
		}
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

	return nil
}

func (a *App) shutdownResources(ctx context.Context) error {
	ctx, span := infraTracing.StartSpan(ctx, "internal:framework:primary:consumer:postgres:run:shutdownResources")
	defer span.End()

	if a.shutdown == nil {
		return nil
	}

	err := a.shutdown()
	if err != nil {
		infraTracing.RecordError(span, err)
	}

	return err
}

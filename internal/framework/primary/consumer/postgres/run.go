package consumer

import (
	infraTracing "codebase-app/internal/infrastructure/tracing"
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
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

	go func() {
		<-ctx.Done()
		logShutdownHandlers(ctx, a.handlers)
		if a.router == nil {
			return
		}
		if closeErr := a.router.Close(); closeErr != nil {
			log.Ctx(ctx).Error().Err(closeErr).Msg("failed to close message router")
		}
	}()

	err = a.router.Run(ctx)
	if err != nil && !errors.Is(err, context.Canceled) {
		infraTracing.RecordError(span, err)
	}

	if shutdownErr := a.shutdownResources(ctx); shutdownErr != nil {
		if err == nil {
			err = shutdownErr
		} else {
			log.Ctx(ctx).Error().Err(shutdownErr).Msg("failed to shutdown resources cleanly")
		}
	}

	return err
}

func logShutdownHandlers(ctx context.Context, handlers []consumerHandlerInfo) {
	if len(handlers) == 0 {
		log.Ctx(ctx).Info().Msg("consumer shutdown requested with no registered handlers")
		return
	}

	log.Ctx(ctx).Info().
		Int("handler_count", len(handlers)).
		Msg("consumer shutdown requested")

	for _, handler := range handlers {
		log.Ctx(ctx).Info().
			Str("handler_name", handler.Name).
			Str("topic", handler.Topic).
			Str("consumer_group", handler.ConsumerGroup).
			Msg("stopping consumer handler")
	}
}

func (a *App) validate() error {
	if a.router == nil {
		return errors.New("message router is required")
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
		if a.routerClose == nil {
			return nil
		}
		return a.routerClose()
	}

	err := a.shutdown()
	if err != nil {
		infraTracing.RecordError(span, err)
	}

	if a.routerClose != nil {
		closeErr := a.routerClose()
		if closeErr != nil {
			infraTracing.RecordError(span, closeErr)
			if err == nil {
				err = closeErr
			} else {
				err = errors.Join(err, closeErr)
			}
		}
	}

	return err
}

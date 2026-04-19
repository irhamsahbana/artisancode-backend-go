package consumer

import (
	infraTracing "codebase-app/internal/infrastructure/tracing"
	"context"
	"os"
	"os/signal"

	"github.com/rs/zerolog/log"
)

func (a *App) waitForShutdown(ctx context.Context) error {
	ctx, span := infraTracing.StartSpan(ctx, "consumer.WaitForShutdown")
	defer span.End()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Ctx(ctx).Info().Msg("Consumer gracefully stopped")

	if a.shutdown == nil {
		return nil
	}

	err := a.shutdown()
	if err != nil {
		infraTracing.RecordError(span, err)
	}

	return err
}

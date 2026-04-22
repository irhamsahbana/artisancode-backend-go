package consumer

import (
	infraTracing "codebase-app/internal/infrastructure/tracing"
	"context"

	"github.com/rs/zerolog/log"
)

func (a *App) waitForShutdown(ctx context.Context) error {
	ctx, span := infraTracing.StartSpan(ctx, "internal:framework:primary:consumer:postgres:shutdown:waitForShutdown")
	defer span.End()

	<-ctx.Done()

	if ctx.Err() == nil {
		return nil
	}

	log.Ctx(ctx).Info().Err(ctx.Err()).Msg("Consumer shutdown signal received")
	return nil
}

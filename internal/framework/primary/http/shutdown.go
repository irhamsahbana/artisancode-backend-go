package http

import (
	"context"

	"github.com/rs/zerolog/log"
)

func (a *App) waitForShutdown(ctx context.Context) error {
	log.Ctx(ctx).Info().Err(ctx.Err()).Msg("HTTP shutdown requested")

	var err error
	if a.shutdown != nil {
		err = a.shutdown()
		if err != nil {
			return err
		}
	}

	log.Ctx(ctx).Info().Msg("Server gracefully stopped")

	return err
}

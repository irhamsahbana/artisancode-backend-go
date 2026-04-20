package http

import (
	"context"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/rs/zerolog/log"
)

func (a *App) waitForShutdown(ctx context.Context) error {
	quit := make(chan os.Signal, 1)

	shutdownSignals := []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGINT}
	if runtime.GOOS == "windows" {
		shutdownSignals = []os.Signal{os.Interrupt}
	}

	signal.Notify(quit, shutdownSignals...)
	<-quit
	log.Ctx(ctx).Info().Msg("Server is shutting down ...")

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

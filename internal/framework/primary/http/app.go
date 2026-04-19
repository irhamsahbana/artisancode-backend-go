package http

import (
	infraTracing "codebase-app/internal/infrastructure/tracing"
	integrationPorts "codebase-app/internal/ports/secondary/integration"
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type App struct {
	bus      integrationPorts.MessagePublisher
	shutdown func() error
}

type AppConfig struct {
	Bus      integrationPorts.MessagePublisher
	Shutdown func() error
}

func NewApp(cfg AppConfig) *App {
	return &App{
		bus:      cfg.Bus,
		shutdown: cfg.Shutdown,
	}
}

func (a *App) Run(ctx context.Context, appName string, appVersion string, appEnvironment string, port string, db *sqlx.DB) error {
	ctx, span := infraTracing.StartSpan(ctx, "http.Run")
	defer span.End()

	app, err := a.build(appName, appVersion, appEnvironment, db)
	if err != nil {
		infraTracing.RecordError(span, err)
		return err
	}
	defer func() {
		if err := a.bus.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close HTTP message bus")
		}
	}()

	go a.serve(app, port)

	err = a.waitForShutdown(ctx)
	if err != nil {
		infraTracing.RecordError(span, err)
	}

	return err
}

func (a *App) serve(app *fiber.App, port string) {
	log.Info().Msgf("Server is running on port %s", port)
	log.Info().Msgf("Connect via: http://localhost:%s", port)
	for _, ip := range getLocalIPv4s() {
		log.Info().Msgf("Connect via: http://%s:%s", ip, port)
	}
	if err := app.Listen(":" + port); err != nil {
		log.Fatal().Msgf("Error while starting server: %v", err)
	}
}

package http

import (
	integrationPorts "codebase-app/internal/ports/integration"
	"context"

	"github.com/gofiber/fiber/v2"
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

func (a *App) Run(
	ctx context.Context,
	appName string,
	appVersion string,
	appEnvironment string,
	port string,
) error {
	app, err := a.build(ctx, appName, appVersion, appEnvironment)
	if err != nil {
		return err
	}

	go a.serve(app, port)

	err = a.waitForShutdown(ctx)
	if err != nil {
		return err
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

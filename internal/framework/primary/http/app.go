package http

import (
	integrationPorts "codebase-app/internal/ports/integration"
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type App struct {
	bus               integrationPorts.MessagePublisher
	shutdown          func() error
	buildFn           func(appName string, appVersion string, appEnvironment string) (*fiber.App, error)
	serveFn           func(app *fiber.App, port string) error
	waitForShutdownFn func(ctx context.Context) error
	listenFn          func(app *fiber.App, addr string) error
	localIPv4sFn      func() []string
}

type AppConfig struct {
	Bus      integrationPorts.MessagePublisher
	Shutdown func() error
}

func NewApp(cfg AppConfig) *App {
	app := &App{
		bus:      cfg.Bus,
		shutdown: cfg.Shutdown,
	}
	app.buildFn = app.build
	app.serveFn = app.serve
	app.waitForShutdownFn = app.waitForShutdown
	app.listenFn = func(app *fiber.App, addr string) error {
		return app.Listen(addr)
	}
	app.localIPv4sFn = getLocalIPv4s

	return app
}

func (a *App) Run(
	ctx context.Context,
	appName string,
	appVersion string,
	appEnvironment string,
	port string,
) error {
	app, err := a.buildFn(appName, appVersion, appEnvironment)
	if err != nil {
		return err
	}

	serveErrCh := make(chan error, 1)
	go func() {
		serveErrCh <- a.serveFn(app, port)
	}()

	shutdownErrCh := make(chan error, 1)
	go func() {
		shutdownErrCh <- a.waitForShutdownFn(ctx)
	}()

	select {
	case err := <-serveErrCh:
		return err
	case err := <-shutdownErrCh:
		return err
	}
}

func (a *App) serve(app *fiber.App, port string) error {
	log.Info().Msgf("Server is running on port %s", port)
	log.Info().Msgf("Connect via: http://localhost:%s", port)
	for _, ip := range a.localIPv4sFn() {
		log.Info().Msgf("Connect via: http://%s:%s", ip, port)
	}

	return a.listenFn(app, ":"+port)
}

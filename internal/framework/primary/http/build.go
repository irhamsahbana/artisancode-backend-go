package http

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/infrastructure"
	infraMetrics "codebase-app/internal/infrastructure/metrics"
	"codebase-app/internal/middleware"
	"codebase-app/internal/setup"
	"codebase-app/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

var syncHTTPAdapters = func(app *fiber.App) {
	adapter.Adapters.Sync(
		adapter.WithRestServer(app),
		adapter.WithPostgres(),
		adapter.WithMessagePublisher(),
		adapter.WithValidator(validator.NewValidator()),
		adapter.WithRestServer(app),
		adapter.WithOpenAISDK(),
		adapter.WithEmailSender(),
		adapter.WithStorage(),
	)
}

var registerHTTPDependencies = setup.HttpDependencies

func (a *App) build(appName string, appVersion string, appEnvironment string) (*fiber.App, error) {
	app := fiber.New()
	syncHTTPAdapters(app)

	metrics := infraMetrics.New(appName, appVersion, appEnvironment)

	app.Use(middleware.CORS())
	app.Use(middleware.RequestID)
	app.Use(middleware.WithAppLogger(log.Logger))
	app.Use(middleware.WithRequestLanguage())
	app.Use(middleware.LocalizeJSONResponse())
	app.Use(middleware.WithHTTPMetrics(metrics))
	app.Use(middleware.WithTracing(appName))
	app.Use(middleware.Recover())
	app.Use(middleware.WithAccessLog(infrastructure.AccessLogger))

	app.Get("/metrics", metrics.Handler())

	registerHTTPDependencies()

	return app, nil
}

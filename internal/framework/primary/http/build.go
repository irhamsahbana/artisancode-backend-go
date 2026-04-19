package http

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/infrastructure"
	infraMetrics "codebase-app/internal/infrastructure/metrics"
	storage "codebase-app/internal/integration/storage"
	"codebase-app/internal/middleware"
	"codebase-app/internal/setup"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (a *App) build(appName string, appVersion string, appEnvironment string, db *sqlx.DB) (*fiber.App, error) {
	app := fiber.New()
	adapter.Adapters.Sync(
		adapter.WithRestServer(app),
		adapter.WithOpenAISDK(),
		adapter.WithStorage(),
	)

	s3 := storage.NewStorageIntegration(adapter.Adapters.Storage)
	metrics := infraMetrics.New(appName, appVersion, appEnvironment)

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,PATCH,OPTIONS,HEAD",
		AllowHeaders: "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin,Authorization",
	}))
	app.Use(middleware.RequestID)
	app.Use(middleware.WithAppLogger(log.Logger))
	app.Use(middleware.WithRequestLanguage())
	app.Use(middleware.LocalizeJSONResponse())
	app.Use(middleware.WithHTTPMetrics(metrics))
	app.Use(middleware.WithTracing(appName))
	app.Use(middleware.Recover())
	app.Use(middleware.WithAccessLog(infrastructure.AccessLogger))

	app.Get("/metrics", metrics.Handler())

	setup.HttpDependencies(
		app,
		db,
		s3,
		a.bus,
	)

	return app, nil
}

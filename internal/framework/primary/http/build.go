package http

import (
	"codebase-app/internal/adapter"
	messagebus "codebase-app/internal/framework/secondary/publisher/messagebus"
	"codebase-app/internal/infrastructure"
	storage "codebase-app/internal/integration/storage"
	"codebase-app/internal/middleware"
	integrationPorts "codebase-app/internal/ports/secondary/integration"
	"codebase-app/internal/setup"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (a *App) build(appName string, appEnvironment string, db *sqlx.DB) (*fiber.App, integrationPorts.MessagePublisher, error) {
	app := fiber.New()
	adapter.Adapters.Sync(
		adapter.WithRestServer(app),
		adapter.WithOpenAISDK(),
		adapter.WithStorage(),
	)

	bus, err := messagebus.NewPublisher(db)
	if err != nil {
		return nil, nil, err
	}

	s3 := storage.NewStorageIntegration(adapter.Adapters.Storage)

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,PATCH,OPTIONS,HEAD",
		AllowHeaders: "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin,Authorization",
	}))
	app.Use(middleware.RequestID)
	app.Use(middleware.WithAppLogger(log.Logger))
	app.Use(middleware.WithRequestLanguage())
	app.Use(middleware.LocalizeJSONResponse())
	app.Use(middleware.WithTracing(appName))
	app.Use(middleware.Recover())
	app.Use(middleware.WithAccessLog(infrastructure.AccessLogger))

	metricTitle := appName + " " + appEnvironment + " " + "Metrics"
	app.Get("/metrics", monitor.New(monitor.Config{Title: metricTitle}))

	setup.HttpDependencies(
		app,
		db,
		s3,
		bus,
	)

	return app, bus, nil
}

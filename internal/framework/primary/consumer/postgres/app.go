package consumer

import (
	corePorts "codebase-app/internal/ports/core"
	integrationPorts "codebase-app/internal/ports/integration"
)

type App struct {
	emailSubscription  integrationPorts.MessageBusSubscription
	exportSubscription integrationPorts.MessageBusSubscription
	exportCore         corePorts.ExportJobCore
	shutdown           func() error
}

type AppConfig struct {
	EmailSubscription  integrationPorts.MessageBusSubscription
	ExportSubscription integrationPorts.MessageBusSubscription
	ExportCore         corePorts.ExportJobCore
	Shutdown           func() error
}

func NewApp(cfg AppConfig) *App {
	return &App{
		emailSubscription:  cfg.EmailSubscription,
		exportSubscription: cfg.ExportSubscription,
		exportCore:         cfg.ExportCore,
		shutdown:           cfg.Shutdown,
	}
}

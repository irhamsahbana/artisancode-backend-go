package consumer

import (
	corePorts "codebase-app/internal/ports/core"
	integrationPorts "codebase-app/internal/ports/secondary/integration"
)

type App struct {
	emailSubscription   integrationPorts.MessageBusSubscription
	exportSubscription  integrationPorts.MessageBusSubscription
	exportCore          corePorts.ExportJobCore
	exportPublisher     integrationPorts.MessagePublisher
	subscriptionManager integrationPorts.MessageSubscriptionManager
	shutdown            func() error
}

type AppConfig struct {
	EmailSubscription   integrationPorts.MessageBusSubscription
	ExportSubscription  integrationPorts.MessageBusSubscription
	ExportCore          corePorts.ExportJobCore
	ExportPublisher     integrationPorts.MessagePublisher
	SubscriptionManager integrationPorts.MessageSubscriptionManager
	Shutdown            func() error
}

func NewApp(cfg AppConfig) *App {
	return &App{
		emailSubscription:   cfg.EmailSubscription,
		exportSubscription:  cfg.ExportSubscription,
		exportCore:          cfg.ExportCore,
		exportPublisher:     cfg.ExportPublisher,
		subscriptionManager: cfg.SubscriptionManager,
		shutdown:            cfg.Shutdown,
	}
}

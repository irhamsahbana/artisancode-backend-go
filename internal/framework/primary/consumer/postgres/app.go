package consumer

import (
	corePorts "codebase-app/internal/ports/core"

	watermillMessage "github.com/ThreeDotsLabs/watermill/message"
)

type App struct {
	router      *watermillMessage.Router
	routerClose func() error
	exportCore  corePorts.ExportJobCore
	handlers    []consumerHandlerInfo
	shutdown    func() error
}

type AppConfig struct {
	Router      *watermillMessage.Router
	RouterClose func() error
	ExportCore  corePorts.ExportJobCore
	Handlers    []consumerHandlerInfo
	Shutdown    func() error
}

type consumerHandlerInfo struct {
	Name          string
	Topic         string
	ConsumerGroup string
	Description   string
}

func NewApp(cfg AppConfig) *App {
	return &App{
		router:      cfg.Router,
		routerClose: cfg.RouterClose,
		exportCore:  cfg.ExportCore,
		handlers:    cfg.Handlers,
		shutdown:    cfg.Shutdown,
	}
}

package exportjob

import (
	"codebase-app/internal/entity/common"
	postgresBus "codebase-app/internal/framework/secondary/publisher/postgres"
	corePorts "codebase-app/internal/ports/core"
)

type Adapter struct {
	core corePorts.ExportJobCore
}

type Config struct {
	Core corePorts.ExportJobCore
}

func NewAdapter(config Config) *Adapter {
	return &Adapter{
		core: config.Core,
	}
}

func (a *Adapter) HandlerDefinitions() []postgresBus.ConsumerHandlerDefinition {
	return []postgresBus.ConsumerHandlerDefinition{
		{
			SubscriptionConfig: common.MessageBusSubscriptionConfig{
				Topics:              []string{common.MessageTopicExportJobRequested},
				ConsumerGroup:       common.MessageConsumerGroupExportJobService,
				ConsumerDescription: "Export job consumer",
			},
			HandlerName: "process_export_job_requested",
			Topic:       common.MessageTopicExportJobRequested,
			Handler:     a.RequestedHandler,
		},
	}
}

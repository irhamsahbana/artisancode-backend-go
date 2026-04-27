package consumer

import (
	"codebase-app/internal/adapter"
	postgresBus "codebase-app/internal/framework/secondary/publisher/postgres"
)

func (a *App) registerConsumerHandlers(
	busConfig postgresBus.Config,
	definitions []postgresBus.ConsumerHandlerDefinition,
) error {
	for _, definition := range definitions {
		if err := postgresBus.AddConsumerHandler(
			a.router,
			adapter.Adapters.Postgres,
			busConfig,
			definition.SubscriptionConfig,
			definition.HandlerName,
			definition.Topic,
			definition.Handler,
		); err != nil {
			return err
		}

		a.handlers = append(a.handlers, consumerHandlerInfo{
			Name:          definition.HandlerName,
			Topic:         definition.Topic,
			ConsumerGroup: definition.SubscriptionConfig.ConsumerGroup,
			Description:   definition.SubscriptionConfig.ConsumerDescription,
		})
	}

	return nil
}

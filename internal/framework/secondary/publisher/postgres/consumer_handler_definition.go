package postgres

import (
	"codebase-app/internal/entity/common"

	watermillMessage "github.com/ThreeDotsLabs/watermill/message"
)

type ConsumerHandlerDefinition struct {
	SubscriptionConfig common.MessageBusSubscriptionConfig
	HandlerName        string
	Topic              string
	Handler            watermillMessage.NoPublishHandlerFunc
}

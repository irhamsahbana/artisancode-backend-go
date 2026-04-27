package userinvitation

import (
	"codebase-app/internal/entity/common"
	postgresBus "codebase-app/internal/framework/secondary/publisher/postgres"
)

type Adapter struct{}

func NewAdapter() *Adapter {
	return &Adapter{}
}

func (a *Adapter) HandlerDefinitions() []postgresBus.ConsumerHandlerDefinition {
	return []postgresBus.ConsumerHandlerDefinition{
		{
			SubscriptionConfig: common.MessageBusSubscriptionConfig{
				Topics:              []string{common.MessageTopicEmailInvitation},
				ConsumerGroup:       common.MessageConsumerGroupNotificationService,
				ConsumerDescription: "Notification service consumer",
			},
			HandlerName: "send_invitation_email",
			Topic:       common.MessageTopicEmailInvitation,
			Handler:     a.InvitationHandler,
		},
	}
}

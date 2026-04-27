package user

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
				Topics:              []string{common.MessageTopicEmailVerification},
				ConsumerGroup:       common.MessageConsumerGroupNotificationService,
				ConsumerDescription: "Notification service consumer",
			},
			HandlerName: "send_email_verification",
			Topic:       common.MessageTopicEmailVerification,
			Handler:     a.EmailVerificationHandler,
		},
		{
			SubscriptionConfig: common.MessageBusSubscriptionConfig{
				Topics:              []string{common.MessageTopicEmailForgotPassword},
				ConsumerGroup:       common.MessageConsumerGroupNotificationService,
				ConsumerDescription: "Notification service consumer",
			},
			HandlerName: "send_forgot_password_email",
			Topic:       common.MessageTopicEmailForgotPassword,
			Handler:     a.ForgotPasswordHandler,
		},
	}
}

package integration

import (
	"context"

	"codebase-app/internal/entity/common"
)

type MessagePublisher interface {
	PublishJSON(ctx context.Context, subject string, payload any) error
	Close() error
}

type MessageSubscriptionManager interface {
	CreateSubscription(ctx context.Context, cfg common.MessageBusSubscriptionConfig) (MessageBusSubscription, error)
	Close() error
}

type MessageBusSubscription interface {
	Consume(ctx context.Context, handler func(context.Context, MessageBusMessage)) (MessageBusSubscriptionContext, error)
}

type MessageBusSubscriptionContext interface {
	Stop()
}

type MessageBusMessage interface {
	Subject() string
	Data() []byte
	Headers() map[string][]string
	Ack(ctx context.Context) error
	Nak(ctx context.Context, reason string) error
}

package integration

import (
	"context"
	"time"
)

type MessagePublisher interface {
	PublishJSON(ctx context.Context, subject string, payload any) error
	Close() error
}

type MessageSubscriptionManager interface {
	CreateSubscription(ctx context.Context, cfg MessageBusSubscriptionConfig) (MessageBusSubscription, error)
	Close() error
}

type MessageBusSubscriptionConfig struct {
	StreamName          string
	StreamDescription   string
	Subjects            []string
	MaxBytes            int64
	MaxAge              time.Duration
	ConsumerName        string
	Durable             string
	ConsumerDescription string
}

type MessageBusSubscription interface {
	Consume(handler func(MessageBusMessage)) (MessageBusSubscriptionContext, error)
}

type MessageBusSubscriptionContext interface {
	Stop()
}

type MessageBusMessage interface {
	Subject() string
	Data() []byte
	Headers() map[string][]string
	Ack() error
	Nak(reason string) error
}

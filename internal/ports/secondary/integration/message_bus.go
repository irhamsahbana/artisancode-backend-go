package integration

import (
	"context"
	"time"
)

type MessageBus interface {
	PublishJSON(ctx context.Context, subject string, payload any) error
	CreateConsumer(ctx context.Context, cfg MessageBusConsumerConfig) (MessageBusConsumer, error)
	Close() error
}

type MessagePublisher interface {
	PublishJSON(ctx context.Context, subject string, payload any) error
	Close() error
}

type MessageConsumerManager interface {
	CreateConsumer(ctx context.Context, cfg MessageBusConsumerConfig) (MessageBusConsumer, error)
	Close() error
}

type MessageBusConsumerConfig struct {
	StreamName          string
	StreamDescription   string
	Subjects            []string
	MaxBytes            int64
	MaxAge              time.Duration
	ConsumerName        string
	Durable             string
	ConsumerDescription string
}

type MessageBusConsumer interface {
	Consume(handler func(MessageBusMessage)) (MessageBusConsumeContext, error)
}

type MessageBusConsumeContext interface {
	Stop()
}

type MessageBusMessage interface {
	Subject() string
	Data() []byte
	Headers() map[string][]string
	Ack() error
	Nak() error
}

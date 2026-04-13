package natsjetstream

import (
	"context"

	"codebase-app/internal/infrastructure/tracing"
	integrationPorts "codebase-app/internal/ports/secondary/integration"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog/log"
)

type consumerManager struct {
	client *client
}

var _ integrationPorts.MessageConsumerManager = &consumerManager{}

func NewConsumerManager(url string) (integrationPorts.MessageConsumerManager, error) {
	client, err := newClient(url)
	if err != nil {
		return nil, err
	}

	return &consumerManager{
		client: client,
	}, nil
}

func (m *consumerManager) CreateConsumer(ctx context.Context, cfg integrationPorts.MessageBusConsumerConfig) (integrationPorts.MessageBusConsumer, error) {
	ctx, span := tracing.StartSpan(ctx, "nats.CreateConsumer")
	defer span.End()

	stream, err := m.getOrCreateStream(ctx, cfg)
	if err != nil {
		return nil, err
	}

	consumer, err := m.createOrUpdateConsumer(ctx, stream, cfg)
	if err != nil {
		return nil, err
	}

	return &natsConsumer{consumer: consumer}, nil
}

func (m *consumerManager) getOrCreateStream(ctx context.Context, cfg integrationPorts.MessageBusConsumerConfig) (jetstream.Stream, error) {
	_, err := m.client.js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:        cfg.StreamName,
		Description: cfg.StreamDescription,
		Subjects:    cfg.Subjects,
		MaxBytes:    cfg.MaxBytes,
		MaxAge:      cfg.MaxAge,
	})
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("Failed to create or update stream")
		return nil, err
	}

	stream, err := m.client.js.Stream(ctx, cfg.StreamName)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("Failed to get stream")
		return nil, err
	}

	return stream, nil
}

func (m *consumerManager) createOrUpdateConsumer(ctx context.Context, stream jetstream.Stream, cfg integrationPorts.MessageBusConsumerConfig) (jetstream.Consumer, error) {
	consumer, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Name:        cfg.ConsumerName,
		Durable:     cfg.Durable,
		Description: cfg.ConsumerDescription,
	})
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("Failed to create or update consumer")
		return nil, err
	}

	return consumer, nil
}

func (m *consumerManager) Close() error {
	return m.client.Close()
}

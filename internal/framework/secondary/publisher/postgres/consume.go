package postgres

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"
	"sync/atomic"

	"github.com/ThreeDotsLabs/watermill"
	watermillSQL "github.com/ThreeDotsLabs/watermill-sql/v4/pkg/sql"
	watermillMessage "github.com/ThreeDotsLabs/watermill/message"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"

	"codebase-app/internal/entity/common"
	integrationPorts "codebase-app/internal/ports/integration"
)

type postgresConsumer struct {
	db         *sqlx.DB
	cfg        Config
	def        common.MessageBusSubscriptionConfig
	subscriber *watermillSQL.Subscriber
}

var _ integrationPorts.MessageBusSubscription = &postgresConsumer{}

func (c *postgresConsumer) Consume(
	ctx context.Context,
	handler func(context.Context, integrationPorts.MessageBusMessage),
) (integrationPorts.MessageBusSubscriptionContext, error) {
	runCtx, cancel := context.WithCancel(context.WithoutCancel(ctx))

	consumeCtx := &postgresConsumeContext{
		cancel:     cancel,
		subscriber: c.subscriber,
		done:       make(chan struct{}),
	}

	if len(c.def.Topics) == 0 {
		cancel()
		return nil, fmt.Errorf("message subscription topics are required")
	}

	for _, topic := range c.def.Topics {
		messages, err := c.subscriber.Subscribe(runCtx, topic)
		if err != nil {
			cancel()
			log.Ctx(ctx).Error().Err(err).Str("topic", topic).Msg("Failed to subscribe to Watermill SQL topic")
			return nil, err
		}

		consumeCtx.wg.Add(1)
		go func(topic string, messages <-chan *watermillMessage.Message) {
			defer consumeCtx.wg.Done()
			c.consumeMessages(runCtx, consumeCtx, topic, messages, handler)
		}(topic, messages)
	}

	go func() {
		consumeCtx.wg.Wait()
		close(consumeCtx.done)
	}()

	go func() {
		select {
		case <-ctx.Done():
			consumeCtx.cancelForStop()
		case <-consumeCtx.done:
		}
	}()

	return consumeCtx, nil
}

func (c *postgresConsumer) consumeMessages(
	ctx context.Context,
	consumeCtx *postgresConsumeContext,
	topic string,
	messages <-chan *watermillMessage.Message,
	handler func(context.Context, integrationPorts.MessageBusMessage),
) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-messages:
			if !ok {
				return
			}
			if consumeCtx.isStopping() {
				_ = msg.Nack()
				return
			}

			consumeCtx.handlerWg.Add(1)
			c.handleMessage(context.WithoutCancel(ctx), topic, handler, msg)
			consumeCtx.handlerWg.Done()

			if consumeCtx.isStopping() {
				return
			}
		}
	}
}

func (c *postgresConsumer) handleMessage(
	ctx context.Context,
	topic string,
	handler func(context.Context, integrationPorts.MessageBusMessage),
	watermillMsg *watermillMessage.Message,
) {
	msg := newPostgresMessage(topic, watermillMsg)
	msgCtx := ctx
	cancel := func() {}
	if c.cfg.ProcessingTimeout > 0 {
		msgCtx, cancel = context.WithTimeout(ctx, c.cfg.ProcessingTimeout)
	}
	defer cancel()

	defer func() {
		if recovered := recover(); recovered != nil {
			log.Ctx(msgCtx).Error().
				Any("panic", recovered).
				Bytes("stack", debug.Stack()).
				Str("subject", msg.Subject()).
				Msg("Message handler panicked")

			reason := fmt.Sprintf("handler panic: %v", recovered)
			if err := msg.Nak(msgCtx, reason); err != nil {
				log.Ctx(msgCtx).Error().Err(err).Str("subject", msg.Subject()).Msg("Failed to negative-acknowledge panicked message")
			}
		}
	}()

	handler(msgCtx, msg)
}

type postgresConsumeContext struct {
	cancel     context.CancelFunc
	subscriber *watermillSQL.Subscriber
	done       chan struct{}
	stopOnce   sync.Once
	stopping   atomic.Bool
	wg         sync.WaitGroup
	handlerWg  sync.WaitGroup
}

var _ integrationPorts.MessageBusSubscriptionContext = &postgresConsumeContext{}

func (c *postgresConsumeContext) Stop() {
	c.stopOnce.Do(func() {
		c.stopping.Store(true)
		c.handlerWg.Wait()
		c.cancelForStop()
		_ = c.subscriber.Close()
	})
	c.wg.Wait()
}

func (c *postgresConsumeContext) cancelForStop() {
	c.stopping.Store(true)
	c.cancel()
}

func (c *postgresConsumeContext) isStopping() bool {
	return c.stopping.Load()
}

func newWatermillSubscriber(
	db watermillSQL.Beginner,
	cfg Config,
	def common.MessageBusSubscriptionConfig,
) (*watermillSQL.Subscriber, error) {
	ackDeadline := cfg.ProcessingTimeout
	if def.ConsumerGroup == "" {
		return nil, fmt.Errorf("message subscription consumer group is required")
	}

	return watermillSQL.NewSubscriber(
		db,
		watermillSQL.SubscriberConfig{
			ConsumerGroup:  def.ConsumerGroup,
			AckDeadline:    &ackDeadline,
			PollInterval:   cfg.PollInterval,
			ResendInterval: cfg.RetryDelay,
			RetryInterval:  cfg.PollInterval,
			SchemaAdapter: watermillSQL.DefaultPostgreSQLSchema{
				SubscribeBatchSize: cfg.BatchSize,
			},
			OffsetsAdapter:   watermillSQL.DefaultPostgreSQLOffsetsAdapter{},
			InitializeSchema: false,
		},
		watermill.NopLogger{},
	)
}

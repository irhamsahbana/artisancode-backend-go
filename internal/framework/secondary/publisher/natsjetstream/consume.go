package natsjetstream

import (
	integrationPorts "codebase-app/internal/ports/secondary/integration"

	"github.com/nats-io/nats.go/jetstream"
)

type natsConsumer struct {
	consumer jetstream.Consumer
}

var _ integrationPorts.MessageBusConsumer = &natsConsumer{}

func (c *natsConsumer) Consume(handler func(integrationPorts.MessageBusMessage)) (integrationPorts.MessageBusConsumeContext, error) {
	ctx, err := c.consumer.Consume(func(msg jetstream.Msg) {
		handler(&natsMessage{msg: msg})
	})
	if err != nil {
		return nil, err
	}

	return &natsConsumeContext{ctx: ctx}, nil
}

type natsConsumeContext struct {
	ctx jetstream.ConsumeContext
}

var _ integrationPorts.MessageBusConsumeContext = &natsConsumeContext{}

func (c *natsConsumeContext) Stop() {
	c.ctx.Stop()
}

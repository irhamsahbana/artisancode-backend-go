package natsjetstream

import (
	"context"
	"encoding/json"

	"codebase-app/internal/infrastructure/tracing"
	integrationPorts "codebase-app/internal/ports/secondary/integration"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
)

type publisher struct {
	client *client
}

var _ integrationPorts.MessagePublisher = &publisher{}

func NewPublisher(url string) (integrationPorts.MessagePublisher, error) {
	client, err := newClient(url)
	if err != nil {
		return nil, err
	}

	return &publisher{
		client: client,
	}, nil
}

func (p *publisher) PublishJSON(ctx context.Context, subject string, payload any) error {
	ctx, span := tracing.StartSpan(ctx, "nats.PublishJSON")
	defer span.End()

	data, err := json.Marshal(payload)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("Failed to marshal payload to JSON")
		return err
	}

	msg := nats.NewMsg(subject)
	msg.Data = data
	otel.GetTextMapPropagator().Inject(ctx, natsHeaderCarrier(msg.Header))

	_, err = p.client.js.PublishMsg(ctx, msg)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("Failed to publish message")
	}

	return err
}

func (p *publisher) Close() error {
	return p.client.Close()
}

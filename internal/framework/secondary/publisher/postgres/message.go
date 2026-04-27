package postgres

import (
	"context"
	"fmt"

	watermillMessage "github.com/ThreeDotsLabs/watermill/message"

	"codebase-app/internal/infrastructure/tracing"
	integrationPorts "codebase-app/internal/ports/integration"
)

type postgresMessage struct {
	msg     *watermillMessage.Message
	subject string
	headers map[string][]string
}

var _ integrationPorts.MessageBusMessage = &postgresMessage{}

func newPostgresMessage(topic string, msg *watermillMessage.Message) *postgresMessage {
	headers := make(map[string][]string, len(msg.Metadata))
	for key, value := range msg.Metadata {
		if key == "topic" {
			continue
		}

		headers[key] = []string{value}
	}

	return &postgresMessage{
		msg:     msg,
		subject: topic,
		headers: headers,
	}
}

func (m *postgresMessage) Subject() string {
	return m.subject
}

func (m *postgresMessage) Data() []byte {
	return m.msg.Payload
}

func (m *postgresMessage) Headers() map[string][]string {
	return m.headers
}

func (m *postgresMessage) Ack(ctx context.Context) error {
	_, span := tracing.StartSpan(ctx, "internal:framework:secondary:publisher:postgres:message:Ack")
	defer span.End()

	if !m.msg.Ack() {
		return fmt.Errorf("watermill message was already negative-acknowledged")
	}

	return nil
}

func (m *postgresMessage) Nak(ctx context.Context, reason string) error {
	_, span := tracing.StartSpan(ctx, "internal:framework:secondary:publisher:postgres:message:Nak")
	defer span.End()

	m.msg.Metadata.Set("last_error", reason)
	if !m.msg.Nack() {
		return fmt.Errorf("watermill message was already acknowledged")
	}

	return nil
}

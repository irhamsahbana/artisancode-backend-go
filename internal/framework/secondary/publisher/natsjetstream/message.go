package natsjetstream

import (
	integrationPorts "codebase-app/internal/ports/secondary/integration"

	"github.com/nats-io/nats.go/jetstream"
)

type natsMessage struct {
	msg jetstream.Msg
}

var _ integrationPorts.MessageBusMessage = &natsMessage{}

func (m *natsMessage) Subject() string {
	return m.msg.Subject()
}

func (m *natsMessage) Data() []byte {
	return m.msg.Data()
}

func (m *natsMessage) Headers() map[string][]string {
	return map[string][]string(m.msg.Headers())
}

func (m *natsMessage) Ack() error {
	return m.msg.Ack()
}

func (m *natsMessage) Nak() error {
	return m.msg.Nak()
}

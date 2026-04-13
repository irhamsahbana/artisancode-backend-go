package natsjetstream

import (
	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel/propagation"
)

type natsHeaderCarrier nats.Header

var _ propagation.TextMapCarrier = natsHeaderCarrier{}

func (c natsHeaderCarrier) Get(key string) string {
	return nats.Header(c).Get(key)
}

func (c natsHeaderCarrier) Set(key, value string) {
	nats.Header(c).Set(key, value)
}

func (c natsHeaderCarrier) Keys() []string {
	keys := make([]string, 0, len(c))
	for key := range c {
		keys = append(keys, key)
	}

	return keys
}
